package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseHeader(t *testing.T) {
	p := postgreDBToCsv{w: nil}
	lineheader := "id                  | codigo_rastreio | cliente nome |      cliente_email      | valor_total |  status   | pago | metadados |           criado_em"

	header, err := p.parseHeader(lineheader)
	if err != nil {
		t.Errorf("parseHeader returned and error %v\n", err)
	}

	if len(header) != 9 {
		t.Errorf("Should have returned header with 9 columns instead returned %v\n", len(header))
	}

	if header[0] != "id" {
		t.Errorf("First column should be id instead got %v\n", header[0])
	}

	if header[1] != "codigo_rastreio" {
		t.Errorf("First column should be codigo_rastreio instead got %v\n", header[1])
	}

	if header[2] != "cliente nome" {
		t.Errorf("First column should be client nome instead got %v with len %v\n", header[2], len(header[2]))
	}
}

func TestParseRow(t *testing.T) {
	var buffer bytes.Buffer
	p := NewPostgreDBToCsv(&buffer)
	content := `
                  id                  | codigo_rastreio | cliente_nome |      cliente_email      | valor_total |  status   | pago | metadados |           criado_em           
--------------------------------------+-----------------+--------------+-------------------------+-------------+-----------+------+-----------+-------------------------------
 72efcb2f-0076-4a38-bf2e-9d3369d69680 | BR000000001PT   | Cliente 1    | cliente1@exemplo.com    |      612.58 | CANCELADO | f    |           | 2026-08-28 09:45:20.157022+00
 3ae4014c-48a2-4f45-81d7-b751971acc0c | BR000000002PT   | Cliente 2    | cliente2@exemplo.com    |      278.18 | ENVIADO   | f    |           | 2026-09-10 14:50:57.647295+00
 6a6b7767-ebc2-442a-8791-071eef42e1d7 | BR000000003PT   | Cliente, 3    | cliente3@exemplo.com    |      347.07 | ENVIADO   | f    |           | 2026-09-08 18:48:29.998346+00
 5db3d120-29db-4bac-8164-2a0820772b68 | BR000000004PT   | Cliente 4    | cliente4@exemplo.com    |      768.57 | PAGO      | t    |           | 2026-08-30 15:03:55.663908+00
(1000 rows)`

	err := p.parseRows(strings.Split(content, "\n")[2:])
	if err != nil {
		t.Errorf("parseHeader returned and error %v\n", err)
	}

	p.w.Flush()
}
