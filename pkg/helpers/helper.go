package helpers

// AddrChecker хэлпер для заполнения адреса
func AddrChecker(addr string) string {
	if addr == "" {
		addr = ":50051"
	} else if addr[0] != ':' { // если пришло "50051"
		addr = ":" + addr
	}
	return addr
}
