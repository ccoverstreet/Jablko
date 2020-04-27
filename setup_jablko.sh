#!/bin/bash

echo "Starting Jablko Smart Home Configuration"

echo "Configuring Jablko SMS Server"
printf "\tEnter Gmail Address: "
read sms_email
printf "\tEnter Gmail Password: "
read sms_password

echo -e "{\n\t\"email\": \"$sms_email\",\n\t\"password\": \"$sms_password\"\n}" > jablko_sms_config.json


