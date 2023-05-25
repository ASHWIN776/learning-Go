#!/bin/bash

go build -o bookings cmd/webApp/*.go && 
./bookings -dbname=bookings -dbuser=postgres -dbpass=pass@3750