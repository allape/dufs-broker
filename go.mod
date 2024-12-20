module github.com/allape/dufs-broker

go 1.23.3

replace github.com/fclairamb/ftpserverlib => github.com/allape/ftpserverlib v0.0.0-20241220101930-5e126b659b7b

require (
	github.com/allape/go-http-vfs v0.0.0-20241220093456-088cf2f65a2b
	github.com/allape/goenv v0.0.0-20241202051618-ce41afb81ebf
	github.com/allape/gogger v0.0.0-20241208090122-dda745ad2428
	github.com/fclairamb/ftpserverlib v0.25.0
	github.com/fclairamb/go-log v0.5.0
	github.com/go-git/go-billy/v5 v5.6.0
	github.com/pkg/sftp v1.13.7
	github.com/spf13/afero v1.11.0
	github.com/willscott/go-nfs v0.0.2
	golang.org/x/crypto v0.31.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/rasky/go-xdr v0.0.0-20170124162913-1a41d1a06c93 // indirect
	github.com/willscott/go-nfs-client v0.0.0-20240104095149-b44639837b00 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)
