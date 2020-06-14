# Jablko Smart Home

## Setting Up

This can be done interactivelly through a startup script ("setup_jablko.sh").

- jablko_sms_config.json
  - This file will contain the gmail account and password for Jablko to send messages through. It is highly advised that this is just a throwaway gmail account as keeping login information in plain-text and allowing unsecure app access in google accounts is not recommended.

        {
          "database_name": "primary",
          "sms_email": "jablkohome@gmail.com",
          "sms_password": "mypassword"
        }

## Starting Jablko

Run the command `./start_jablko.sh`. Make sure to make the bash script executable beforehand.
