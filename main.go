package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)


func getServerIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "Nie udało się pobrać adresu IP"
	}

	for _, iface := range interfaces {
		// Ignorujemy interfejsy, które są w trybie loopback (localhost)
		if iface.Flags&net.FlagUp == 0 || strings.Contains(iface.Name, "lo") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			// Pomijamy adresy IPv6
			if strings.Contains(addr.String(), ":") {
				continue
			}
			return addr.String()
		}
	}

	return "Brak dostępnych adresów IP"
}

// Funkcja do uzyskania nazwy hosta
func getHostName() string {
	host, err := os.Hostname()
	if err != nil {
		return "Nie udało się pobrać nazwy hosta"
	}
	return host
}

func getAppVersion() string {
	// Pobieramy wersję z zmiennej środowiskowej
	version := os.Getenv("APP_VERSION")

	// Jeśli zmienna środowiskowa nie została ustawiona, zwrócimy domyślną wersję
	if version == "" {
		return "Nie ustawiono wersji"
	}

	return version
}

// Obsługa zapytania HTTP
func handler(w http.ResponseWriter, r *http.Request) {
	ip := getServerIP()
	host := getHostName()
	version := getAppVersion()

    // Ustawiamy odpowiedni nagłówek, aby wskazać, że odpowiedź będzie w formacie HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Wysyłamy odpowiedź w formacie HTML w jednym Fprintf
	fmt.Fprintf(w, `<!DOCTYPE html>
    <html lang="pl">
        <body>
            <p><strong>Adres IP serwera:</strong> %s</p>
            <p><strong>Nazwa serwera:</strong> %s</p>
            <p><strong>Wersja aplikacji:</strong> %s</p>
        </body>
    </html>`, ip, host, version)
}

func main() {
	// Definiujemy endpoint na '/' i przypisujemy do niego funkcję handler
	http.HandleFunc("/", handler)

	// Uruchamiamy serwer na porcie 8080
	fmt.Println("Serwer działa na porcie 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Błąd uruchamiania serwera:", err)
	}
}
