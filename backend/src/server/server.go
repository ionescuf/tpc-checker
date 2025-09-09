package server

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"errors"
	"go2-gkes-tpc/src/config"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	AuthTypeToken       = 1
	AuthTypeCredentials = 2
	AuthTypeLoginConfig = 3
)

type Server struct {
	mux  *http.ServeMux
	conf *config.Config
}

func Main(conf *config.Config) {
	mux := http.NewServeMux()
	s := &Server{mux: mux, conf: conf}
	s.mux.HandleFunc("/upload-to-storage", s.uploadToStorage)
	err := http.ListenAndServe(s.conf.Server.Host, s.mux)
	if err != nil {
		return
	}
}

type JsonResponse struct {
	Success bool         `json:"success"`
	Err     string       `json:"error"`
	Logs    []*LogRecord `json:"logs"`
}

func (s *Server) uploadToStorage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, http.StatusMethodNotAllowed, make([]*LogRecord, 0), errors.New("method not allowed"))
		return
	}

	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("storageFile")
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
		return
	}
	defer file.Close()
	if header.Size > 1024*1024 {
		jsonResponse(w, http.StatusRequestEntityTooLarge, make([]*LogRecord, 0), errors.New("storage file too large"))
		return
	}

	storageOptions := make([]option.ClientOption, 0)

	if len(r.MultipartForm.Value["universeDomain"]) > 0 {
		storageOptions = append(storageOptions, option.WithUniverseDomain(r.MultipartForm.Value["universeDomain"][0]))
	}

	at := AuthTypeToken
	if len(r.MultipartForm.Value["authType"]) > 0 {
		at, err = strconv.Atoi(r.MultipartForm.Value["authType"][0])
		if err != nil {
			jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
			return
		}
	}

	c := NewHttpClient()
	//storageOptions = append(storageOptions, option.WithHTTPClient(c.Client()))

	//storageOptions = append(storageOptions, option.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))))

	if at == AuthTypeToken {
		var ts *TokenSource
		if len(r.MultipartForm.Value["authToken"]) > 0 && r.MultipartForm.Value["authToken"][0] != "" {
			ts = &TokenSource{token: strings.TrimSpace(r.MultipartForm.Value["authToken"][0])}
		} else {
			ts = &TokenSource{token: strings.TrimSpace(s.conf.GCP.Token)}
		}
		storageOptions = append(storageOptions, option.WithTokenSource(ts))
		//tk, err := ts.Token()
		//if err != nil {
		//	jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
		//	return
		//}
		//c.SetToken(tk)
	}

	if at == AuthTypeCredentials {
		credFile, credFileHeader, err := r.FormFile("credFile")
		if err != nil {
			jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
			return
		}
		defer credFile.Close()
		if credFileHeader.Size > 1024*1024 {
			jsonResponse(w, http.StatusRequestEntityTooLarge, make([]*LogRecord, 0), errors.New("credentials file too large"))
			return
		}
		jsonCredData, err := io.ReadAll(credFile)
		if err != nil {
			jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
			return
		}
		//creds, err := google.CredentialsFromJSON(r.Context(), jsonCredData, "https://www.googleapis.com/auth/cloud-platform")
		//if err != nil {
		//	jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
		//	return
		//}

		storageOptions = append(storageOptions, option.WithCredentialsJSON(jsonCredData))

		//tk, err := creds.TokenSource.Token()
		//if err != nil {
		//	jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
		//	return
		//}
		//c.SetToken(tk)
	}

	if at == AuthTypeLoginConfig {
		lcFile, lcFileHeader, err := r.FormFile("loginConfigFile")
		if err != nil {
			jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
			return
		}
		defer lcFile.Close()
		if lcFileHeader.Size > 1024*1024 {
			jsonResponse(w, http.StatusRequestEntityTooLarge, make([]*LogRecord, 0), errors.New("login config file too large"))
			return
		}
		jsonCredData, err := io.ReadAll(lcFile)
		if err != nil {
			jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
			return
		}

		var raw map[string]interface{}
		if err := json.Unmarshal(jsonCredData, &raw); err != nil {
			jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
			return
		}

		if cs, ok := raw["credential_source"]; ok {
			v := cs.(map[string]interface{})
			if _, ok := v["file"]; !ok {
				jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
				return
			}

			tokenFile, tokenFileHeader, err := r.FormFile("tokenFile")
			if err != nil {
				jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
				return
			}
			defer tokenFile.Close()
			if tokenFileHeader.Size > 1024*1024 {
				jsonResponse(w, http.StatusRequestEntityTooLarge, make([]*LogRecord, 0), errors.New("token file too large"))
				return
			}
			ext := filepath.Ext(tokenFileHeader.Filename)

			fname := "./token-file" + ext
			v["file"] = fname
			raw["credential_source"] = v
			updatedCredData, err := json.Marshal(raw)
			if err != nil {
				jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
				return
			}

			os.WriteFile("./cred-data.json", updatedCredData, 0664)

			tokenFileData, err := io.ReadAll(tokenFile)
			if err != nil {
				jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
				return
			}
			err = os.WriteFile(fname, tokenFileData, 0664)
			if err != nil {
				jsonResponse(w, http.StatusInternalServerError, make([]*LogRecord, 0), err)
				return
			}

			storageOptions = append(storageOptions, option.WithCredentialsJSON(updatedCredData))
		} else {
			jsonResponse(w, http.StatusBadRequest, make([]*LogRecord, 0), err)
			return
		}
	}

	var bucketName string
	if len(r.MultipartForm.Value["bucketName"]) > 0 && r.MultipartForm.Value["bucketName"][0] != "" {
		bucketName = strings.TrimSpace(r.MultipartForm.Value["bucketName"][0])
	} else {
		bucketName = strings.TrimSpace(s.conf.GCP.BucketName)
	}

	client, err := storage.NewClient(context.Background(), storageOptions...)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, c.Logs(), err)
		return
	}
	defer client.Close()

	wr := client.Bucket(bucketName).Object(header.Filename).NewWriter(r.Context())

	if _, err := io.Copy(wr, file); err != nil {
		jsonResponse(w, http.StatusInternalServerError, c.Logs(), err)
		return
	}
	if err := wr.Close(); err != nil {
		jsonResponse(w, http.StatusInternalServerError, c.Logs(), err)
		return
	}

	jsonResponse(w, http.StatusOK, c.Logs(), nil)
}

func jsonResponse(w http.ResponseWriter, code int, logs []*LogRecord, e error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	resp := &JsonResponse{Success: true}
	if len(logs) > 0 {
		resp.Logs = logs
	}
	if e != nil {
		resp.Err = e.Error()
		resp.Success = false
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type TokenSource struct {
	token string
}

func (ts *TokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: ts.token,
	}, nil
}
