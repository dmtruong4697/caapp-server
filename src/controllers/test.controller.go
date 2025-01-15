package controllers

import (
	"caapp-server/src/utils/helper"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"

	"crypto/md5"
	"encoding/hex"
)

func TestSendMail(c *gin.Context) {
	mailer := gomail.NewMessage()

	// Configure the sender and recipient
	mailer.SetHeader("From", os.Getenv("MAIL_USERNAME"))
	mailer.SetHeader("To", "truonggduonggmadridista@gmail.com")
	mailer.SetHeader("Subject", "Welcome to CAApp!")

	// Email body in HTML format
	emailBody := `
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body {
					font-family: Arial, sans-serif;
					line-height: 1.6;
					color: #333;
					background-color: #f9f9f9;
					padding: 20px;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					background: #ffffff;
					padding: 20px;
					border-radius: 8px;
					box-shadow: 0 2px 4px rgba(0,0,0,0.1);
				}
				.header {
					text-align: center;
					margin-bottom: 20px;
				}
				.header img {
					max-width: 100px;
				}
				.content {
					font-size: 16px;
				}
				.footer {
					margin-top: 20px;
					font-size: 14px;
					color: #777;
					text-align: center;
				}
				.button {
					display: inline-block;
					padding: 10px 20px;
					color: #fff;
					background-color:rgb(0, 177, 106);
					text-decoration: none;
					border-radius: 5px;
					font-weight: bold;
					margin-top: 10px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<img src="https://img-s-msn-com.akamaized.net/tenant/amp/entityid/AA1w7ZFF.img?w=730&h=427&m=6" alt="CAApp Logo">
					<h1>Welcome to CAApp!</h1>
				</div>
				<div class="content">
					<p>Hello,</p>
					<p>Thank you for signing up with <strong>CAApp</strong>.</p>
					<p>We are thrilled to have you on board!</p>
					<p>Click the button below to explore our services:</p>
					<a class="button" href="https://your-app-url.com">Explore Now</a>
				</div>
				<div class="footer">
					<p>CAApp - Optimizing your user experience.</p>
					<p>Contact us: <a href="mailto:support@caapp.com">support@caapp.com</a></p>
				</div>
			</div>
		</body>
		</html>
	`

	// Set the email body with HTML content
	mailer.SetBody("text/html", emailBody)

	// Configure the SMTP connection
	dialer := gomail.NewDialer(
		os.Getenv("MAIL_HOST"),
		int(helper.StringToUInt(os.Getenv("MAIL_PORT"))),
		os.Getenv("MAIL_USERNAME"),
		os.Getenv("MAIL_PASSWORD"),
	)

	// Send the email
	err := dialer.DialAndSend(mailer)
	if err != nil {
		fmt.Println("Error sending email:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "email_send_failed"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully!"})
}

func TestDiceBear(c *gin.Context) {
	normalizedEmail := url.QueryEscape("duongminhtruong2002.lequydon@gmail.com")
	hasher := md5.New()
	hasher.Write([]byte(normalizedEmail))
	emailHash := hex.EncodeToString(hasher.Sum(nil))

	dicebearURL := fmt.Sprintf("https://api.dicebear.com/5.x/identicon/svg?seed=%s", emailHash)

	c.JSON(http.StatusOK, gin.H{
		"avatar_url": dicebearURL,
	})
}
