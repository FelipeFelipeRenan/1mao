package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

func sendResetPasswordEmail(to, token string) error {
	from := os.Getenv("EMAIL_SERVICE")
	password := os.Getenv("EMAIL_PASSWORD")

	if from == "" || password == "" {
		return fmt.Errorf("⚠️ EMAIL_SERVICE ou EMAIL_PASSWORD não estão definidos")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Criando a mensagem
	subject := "Subject: 🔑 Redefinição de Senha\n"
	mime := "MIME-Version: 1.0\nContent-Type: text/plain; charset=\"utf-8\"\n\n"
	body := fmt.Sprintf("Olá,\n\nRecebemos um pedido para redefinir sua senha.\nToken: %s\n\nSe não foi você, ignore este e-mail.\n\nEquipe 1Mão", token)
	message := []byte(subject + mime + body)

	// Conectando ao servidor SMTP
	auth := smtp.PlainAuth("", from, password, smtpHost)
	tlsConfig := &tls.Config{ServerName: smtpHost}

	conn, err := smtp.Dial(smtpHost + ":" + smtpPort)
	if err != nil {
		return fmt.Errorf("❌ Erro ao conectar ao servidor SMTP: %v", err)
	}
	defer conn.Close()

	// Iniciando comunicação segura e autenticação
	if err = conn.StartTLS(tlsConfig); err != nil {
		return err
	}
	if err = conn.Auth(auth); err != nil {
		return err
	}

	// Enviando e-mail
	if err = conn.Mail(from); err != nil {
		return err
	}
	if err = conn.Rcpt(to); err != nil {
		return err
	}

	w, err := conn.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write(message)
	if err != nil {
		return err
	}

	fmt.Println("✅ E-mail enviado com sucesso para:", to)
	return conn.Quit()
}
