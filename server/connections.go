package server

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func handle_conn(s *Server, conn *net.Conn, dir string) {
	defer func() {
		if err := s.Close(*conn); err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}()
	for {
		buf := make([]byte, 1024)
		request, err := s.Read(*conn, buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Connection closed by client")
				return
			}
			fmt.Println("Error reading from connection:", err)
			return
		}
		fmt.Println("Received request:\n\n", request)

		lines := strings.Split(request, "\r\n")
		if len(lines) < 1 {
			fmt.Println("Invalid request format")
			return
		}
		request_headers := make(map[string]string)
		i := 0
		line := ""
		for i, line = range lines[1:] {
			if line == "" {
				break
			}
			parts := strings.SplitN(line, ": ", 2)
			if len(parts) != 2 {
				fmt.Println("Invalid header format:", line)
				return
			}
			request_headers[parts[0]] = parts[1]
		}
		body := lines[i+1:]
		start_line := strings.Split(lines[0], " ")
		method := start_line[0]
		path := start_line[1]

		path_parts := strings.Split(path, "/")
		for i, part := range path_parts {
			path_parts[i], err = url.PathUnescape(part)
			if err != nil {
				fmt.Println("Error unescaping path part:", err)
				return
			}
		}

		var status_code int
		var status_text string
		response_headers := make(map[string]string)
		var should_compress bool
		if compression, ok := request_headers["Accept-Encoding"]; ok {
			if strings.Contains(compression, "gzip") {
				should_compress = true
				response_headers["Content-Encoding"] = "gzip"
			}
		}

		response_body := ""
		switch path_parts[1] {
		case "":
			status_code = 200
			status_text = "OK"
			response_headers["Content-Type"] = "text/plain"
		case "echo":
			status_code = 200
			status_text = "OK"
			response_headers["Content-Type"] = "text/plain"
			response_body = path_parts[2]
		case "user-agent":
			status_code = 200
			status_text = "OK"
			if userAgent, ok := request_headers["User-Agent"]; ok {
				response_headers["Content-Type"] = "text/plain"
				response_body = userAgent
			} else {
				status_code = 400
				status_text = "Bad Request"
				response_headers["Content-Type"] = "text/plain"
				response_body = "User-Agent header missing"
			}
		case "files":
			switch method {
			case "GET":
				file, err := os.ReadFile(dir + "/" + filepath.Clean(path_parts[2]))
				if err != nil {
					status_code = 404
					status_text = "Not Found"
					response_headers["Content-Type"] = "text/plain"
					response_body = "404 Not Found"
				} else {
					status_code = 200
					status_text = "OK"
					response_headers["Content-Type"] = "application/octet-stream"
					response_body = string(file)
				}
			case "POST":
				cleanPath := filepath.Clean(filepath.Join(path_parts[2:]...))
				fullPath := filepath.Join(dir, cleanPath)
				rel, err := filepath.Rel(dir, fullPath)
				if err != nil || strings.HasPrefix(rel, "..") {
					status_code = 400
					status_text = "Bad Request"
					response_headers["Content-Type"] = "text/plain"
					response_body = "400 Bad Request"
					break
				}
				if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
					status_code = 500
					status_text = "Internal Server Error"
					response_headers["Content-Type"] = "text/plain"
					response_body = "500 Internal Server Error"
					break
				}
				file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					status_code = 500
					status_text = "Internal Server Error"
					response_headers["Content-Type"] = "text/plain"
					response_body = "500 Internal Server Error"
				} else {
					defer file.Close()
					if _, err := file.WriteString(strings.Join(body, "\n")); err != nil {
						status_code = 500
						status_text = "Internal Server Error"
						response_headers["Content-Type"] = "text/plain"
						response_body = "500 Internal Server Error"
					} else {
						status_code = 201
						status_text = "Created"
						response_headers["Content-Type"] = "text/plain"
						response_body = "Created"
					}
				}
			}
		default:
			status_code = 404
			status_text = "Not Found"
			response_headers["Content-Type"] = "text/plain"
			response_body = "404 Not Found"
		}

		// Prepare response body (compress if needed)
		var final_body []byte
		if should_compress && response_body != "" {
			var b bytes.Buffer
			gzWriter := gzip.NewWriter(&b)
			if _, err := gzWriter.Write([]byte(response_body)); err != nil {
				fmt.Println("Error compressing response:", err)
				return
			}
			if err := gzWriter.Close(); err != nil {
				fmt.Println("Error closing gzip writer:", err)
				return
			}
			final_body = b.Bytes()
		} else {
			final_body = []byte(response_body)
		}

		closeConn := strings.ToLower(request_headers["Connection"]) == "close"
		if closeConn {
			response_headers["Connection"] = "close"
		} else {
			response_headers["Connection"] = "keep-alive"
		}

		// Set content length
		response_headers["Content-Length"] = fmt.Sprint(len(final_body))

		// Build response
		response := ""
		response += fmt.Sprintf("HTTP/1.1 %d %s\r\n", status_code, status_text)
		for key, value := range response_headers {
			response += fmt.Sprintf("%s: %s\r\n", key, value)
		}
		response += "\r\n"

		fmt.Print("Sending response:\n\n", response)

		// Send headers
		if _, err = s.Write(*conn, []byte(response)); err != nil {
			fmt.Println("Error writing headers to connection:", err)
			return
		}

		// Send body if present
		if len(final_body) > 0 {
			if _, err = s.Write(*conn, final_body); err != nil {
				fmt.Println("Error writing body to connection:", err)
				return
			}
		}
		if closeConn {
			return
		}
	}
}
