# Usa una immagine di base con Go preinstallato
FROM golang:latest

# Imposta la directory di lavoro nel percorso del codice Go
WORKDIR /go/src/service_registry

# Copia il codice sorgente del servizio Go RPC nella directory di lavoro del container
COPY . .

# Compila il servizio Go RPC e il file main
RUN go build -o registry ./registry


