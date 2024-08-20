#!/bin/bash

# URL da chiamare
url="http://localhost:3000/push_string"

# Numero di thread che vuoi usare
num_threads=10

# Funzione per generare una stringa casuale di lunghezza specificata
generate_random_string() {
  length=$1
  LC_CTYPE=C tr -dc 'a-zA-Z0-9' </dev/urandom | head -c "$length"
}

# Creare uno script temporaneo con la funzione send_request
temp_script=$(mktemp)

cat << 'EOF' > "$temp_script"
#!/bin/bash

# URL da chiamare
url="http://localhost:3000/push_string"

# Funzione per generare una stringa casuale di lunghezza specificata
generate_random_string() {
  length=$1
  LC_CTYPE=C tr -dc 'a-zA-Z0-9' </dev/urandom | head -c "$length"
}

# Funzione che esegue una singola richiesta
send_request() {
  i=$1
  key="key$i"
  value=$(generate_random_string 10)
  json_body=$(printf '{"key": "%s", "value": "%s"}' "$key" "$value")
  curl -s -X POST -H "Content-Type: application/json" -d "$json_body" "$url"
}

# Chiamare la funzione con l'argomento passato
send_request "$1"
EOF

# Rendere lo script temporaneo eseguibile
chmod +x "$temp_script"

# Generare una sequenza di numeri da 1 a 100000 e inviarli in parallelo usando xargs
seq 1 100000 | xargs -n 1 -P $num_threads "$temp_script"

# Rimuovere lo script temporaneo
rm "$temp_script"
