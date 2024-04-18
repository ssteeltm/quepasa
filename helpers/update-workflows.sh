#!/bin/bash

DBUSER=n8nuser
DBPASS=dfdsj23j42345
DBHOST=localhost
DBNAME=n8n_production

psql postgresql://$DBUSER:$DBPASS@$DBHOST/$DBNAME -c "SELECT id FROM public.user LIMIT 1" -tA

cd /root/.n8n
if [ -z $1 ]; then 
	if !(/usr/bin/n8n import:workflow --input=/opt/quepasa-source/extra/n8n+chatwoot/ --separate); then
		exit 1;
	fi	
else 
	if !(/usr/bin/n8n import:workflow --input=/opt/quepasa-source/extra/n8n+chatwoot/ --separate --userId=$1); then
		exit 1;
	fi
fi

echo ""
echo "########################################"
echo "workflows imported with success"

/usr/bin/n8n update:workflow --id 1008 --active=true &>/dev/null
/usr/bin/n8n update:workflow --id 1009 --active=true &>/dev/null
/usr/bin/n8n update:workflow --id 1010 --active=true &>/dev/null
/usr/bin/n8n update:workflow --id 1011 --active=true &>/dev/null
echo "workflows activated with success"

echo ""
echo "*don't forget to open postgres nodes and update the current credentials."

exit 0