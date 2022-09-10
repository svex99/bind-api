$ORIGIN test.com.
$TTL 12d

@ IN SOA ns1 admin ( 1 23 456 78 90 )

@ IN NS ns1 ; this is an ns record

ns1 IN A 10.10.10.10

@ IN MX 100 email

;
; multiline comment
; comment IN A 1.1.1.1
;

@ IN TXT "var=hello world1"

email IN A 20.20.20.20

@ IN NS ns2

@ IN TXT "var=hello world2"
