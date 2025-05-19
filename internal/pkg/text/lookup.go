package text

import (
	"regexp"
)

// Source: https://github.com/s0md3v/Bolt/blob/master/db/hashes.json
var db = []struct {
	regex string
	names []string
}{
	{
		regex: "^[a-f0-9]{4}$",
		names: []string{
			"CRC-16",
			"CRC-16-CCITT",
			"FCS-16",
		},
	},
	{
		regex: "^[a-f0-9]{8}$",
		names: []string{
			"Adler-32",
			"CRC-32B",
			"FCS-32",
			"GHash-32-3",
			"GHash-32-5",
			"FNV-132",
			"Fletcher-32",
			"Joaat",
			"ELF-32",
			"XOR-32",
		},
	},
	{
		regex: "^[a-f0-9]{6}$",
		names: []string{
			"CRC-24",
		},
	},
	{
		regex: "^(\\$crc32\\$[a-f0-9]{8}.)?[a-f0-9]{8}$",
		names: []string{
			"CRC-32",
		},
	},
	{
		regex: "^\\+[a-z0-9\\/.]{12}$",
		names: []string{
			"Eggdrop IRC Bot",
		},
	},
	{
		regex: "^[a-z0-9\\/.]{13}$",
		names: []string{
			"DES(Unix)",
			"Traditional DES",
			"DEScrypt",
		},
	},
	{
		regex: "^[a-f0-9]{16}$",
		names: []string{
			"MySQL323",
			"DES(Oracle)",
			"Half MD5",
			"Oracle 7-10g",
			"FNV-164",
			"CRC-64",
		},
	},
	{
		regex: "^[a-z0-9\\/.]{16}$",
		names: []string{
			"Cisco-PIX(MD5)",
		},
	},
	{
		regex: "^\\([a-z0-9\\/+]{20}\\)$",
		names: []string{
			"Lotus Notes/Domino 6",
		},
	},
	{
		regex: "^_[a-z0-9\\/.]{19}$",
		names: []string{
			"BSDi Crypt",
		},
	},
	{
		regex: "^[a-f0-9]{24}$",
		names: []string{
			"CRC-96(ZIP)",
		},
	},
	{
		regex: "^[a-z0-9\\/.]{24}$",
		names: []string{
			"Crypt16",
		},
	},
	{
		regex: "^(\\$md2\\$)?[a-f0-9]{32}$",
		names: []string{
			"MD2",
		},
	},
	{
		regex: "^[a-f0-9]{32}(:.+)?$",
		names: []string{
			"MD5",
			"MD4",
			"Double MD5",
			"LM",
			"RIPEMD-128",
			"Haval-128",
			"Tiger-128",
			"Skein-256(128)",
			"Skein-512(128)",
			"Lotus Notes/Domino 5",
			"Skype",
			"ZipMonster",
			"PrestaShop",
			"md5(md5(md5($pass)))",
			"md5(strtoupper(md5($pass)))",
			"md5(sha1($pass))",
			"md5($pass.$salt)",
			"md5($salt.$pass)",
			"md5(unicode($pass).$salt)",
			"md5($salt.unicode($pass))",
			"HMAC-MD5 (key = $pass)",
			"HMAC-MD5 (key = $salt)",
			"md5(md5($salt).$pass)",
			"md5($salt.md5($pass))",
			"md5($pass.md5($salt))",
			"md5($salt.$pass.$salt)",
			"md5(md5($pass).md5($salt))",
			"md5($salt.md5($salt.$pass))",
			"md5($salt.md5($pass.$salt))",
			"md5($username.0.$pass)",
		},
	},
	{
		regex: "^(\\$snefru\\$)?[a-f0-9]{32}$",
		names: []string{
			"Snefru-128",
		},
	},
	{
		regex: "^(\\$NT\\$)?[a-f0-9]{32}$",
		names: []string{
			"NTLM",
		},
	},
	{
		regex: "^([^\\\\\\/:*?\"<>|]{1,20}:)?[a-f0-9]{32}(:[^\\\\\\/:*?\"<>|]{1,20})?$",
		names: []string{
			"Domain Cached Credentials",
		},
	},
	{
		regex: "^([^\\\\\\/:*?\"<>|]{1,20}:)?(\\$DCC2\\$10240#[^\\\\\\/:*?\"<>|]{1,20}#)?[a-f0-9]{32}$",
		names: []string{
			"Domain Cached Credentials 2",
		},
	},
	{
		regex: "^{SHA}[a-z0-9\\/+]{27}=$",
		names: []string{
			"SHA-1(Base64)",
			"Netscape LDAP SHA",
		},
	},
	{
		regex: "^\\$1\\$[a-z0-9\\/.]{0,8}\\$[a-z0-9\\/.]{22}(:.*)?$",
		names: []string{
			"MD5 Crypt",
			"Cisco-IOS(MD5)",
			"FreeBSD MD5",
		},
	},
	{
		regex: "^0x[a-f0-9]{32}$",
		names: []string{
			"Lineage II C4",
		},
	},
	{
		regex: "^\\$H\\$[a-z0-9\\/.]{31}$",
		names: []string{
			"phpBB v3.x",
			"Wordpress v2.6.0/2.6.1",
			"PHPass' Portable Hash",
		},
	},
	{
		regex: "^\\$P\\$[a-z0-9\\/.]{31}$",
		names: []string{
			"Wordpress \u2265 v2.6.2",
			"Joomla \u2265 v2.5.18",
			"PHPass' Portable Hash",
		},
	},
	{
		regex: "^[a-f0-9]{32}:[a-z0-9]{2}$",
		names: []string{
			"osCommerce",
			"xt:Commerce",
		},
	},
	{
		regex: "^\\$apr1\\$[a-z0-9\\/.]{0,8}\\$[a-z0-9\\/.]{22}$",
		names: []string{
			"MD5(APR)",
			"Apache MD5",
			"md5apr1",
		},
	},
	{
		regex: "^{smd5}[a-z0-9$\\/.]{31}$",
		names: []string{
			"AIX(smd5)",
		},
	},
	{
		regex: "^[a-f0-9]{32}:[a-f0-9]{32}$",
		names: []string{
			"WebEdition CMS",
		},
	},
	{
		regex: "^[a-f0-9]{32}:.{5}$",
		names: []string{
			"IP.Board \u2265 v2+",
		},
	},
	{
		regex: "^[a-f0-9]{32}:.{8}$",
		names: []string{
			"MyBB \u2265 v1.2+",
		},
	},
	{
		regex: "^[a-z0-9]{34}$",
		names: []string{
			"CryptoCurrency(Address)",
		},
	},
	{
		regex: "^[a-f0-9]{40}(:.+)?$",
		names: []string{
			"SHA-1",
			"Double SHA-1",
			"RIPEMD-160",
			"Haval-160",
			"Tiger-160",
			"HAS-160",
			"LinkedIn",
			"Skein-256(160)",
			"Skein-512(160)",
			"MangosWeb Enhanced CMS",
			"sha1(sha1(sha1($pass)))",
			"sha1(md5($pass))",
			"sha1($pass.$salt)",
			"sha1($salt.$pass)",
			"sha1(unicode($pass).$salt)",
			"sha1($salt.unicode($pass))",
			"HMAC-SHA1 (key = $pass)",
			"HMAC-SHA1 (key = $salt)",
			"sha1($salt.$pass.$salt)",
		},
	},
	{
		regex: "^\\*[a-f0-9]{40}$",
		names: []string{
			"MySQL5.x",
			"MySQL4.1",
		},
	},
	{
		regex: "^[a-z0-9]{43}$",
		names: []string{
			"Cisco-IOS(SHA-256)",
		},
	},
	{
		regex: "^{SSHA}[a-z0-9\\/+]{38}==$",
		names: []string{
			"SSHA-1(Base64)",
			"Netscape LDAP SSHA",
			"nsldaps",
		},
	},
	{
		regex: "^[a-z0-9=]{47}$",
		names: []string{
			"Fortigate(FortiOS)",
		},
	},
	{
		regex: "^[a-f0-9]{48}$",
		names: []string{
			"Haval-192",
			"Tiger-192",
			"SHA-1(Oracle)",
			"OSX v10.4",
			"OSX v10.5",
			"OSX v10.6",
		},
	},
	{
		regex: "^[a-f0-9]{51}$",
		names: []string{
			"Palshop CMS",
		},
	},
	{
		regex: "^[a-z0-9]{51}$",
		names: []string{
			"CryptoCurrency(PrivateKey)",
		},
	},
	{
		regex: "^{ssha1}[0-9]{2}\\$[a-z0-9$\\/.]{44}$",
		names: []string{
			"AIX(ssha1)",
		},
	},
	{
		regex: "^0x0100[a-f0-9]{48}$",
		names: []string{
			"MSSQL(2005)",
			"MSSQL(2008)",
		},
	},
	{
		regex: "^(\\$md5,rounds=[0-9]+\\$|\\$md5\\$rounds=[0-9]+\\$|\\$md5\\$)[a-z0-9\\/.]{0,16}(\\$|\\$\\$)[a-z0-9\\/.]{22}$",
		names: []string{
			"Sun MD5 Crypt",
		},
	},
	{
		regex: "^[a-f0-9]{56}$",
		names: []string{
			"SHA-224",
			"Haval-224",
			"SHA3-224",
			"Skein-256(224)",
			"Skein-512(224)",
		},
	},
	{
		regex: "^(\\$2[axy]|\\$2)\\$[0-9]{2}\\$[a-z0-9\\/.]{53}$",
		names: []string{
			"Blowfish(OpenBSD)",
			"Woltlab Burning Board 4.x",
			"bcrypt",
		},
	},
	{
		regex: "^[a-f0-9]{40}:[a-f0-9]{16}$",
		names: []string{
			"Android PIN",
		},
	},
	{
		regex: "^(S:)?[a-f0-9]{40}(:)?[a-f0-9]{20}$",
		names: []string{
			"Oracle 11g/12c",
		},
	},
	{
		regex: "^\\$bcrypt-sha256\\$(2[axy]|2)\\,[0-9]+\\$[a-z0-9\\/.]{22}\\$[a-z0-9\\/.]{31}$",
		names: []string{
			"bcrypt(SHA-256)",
		},
	},
	{
		regex: "^[a-f0-9]{32}:.{3}$",
		names: []string{
			"vBulletin < v3.8.5",
		},
	},
	{
		regex: "^[a-f0-9]{32}:.{30}$",
		names: []string{
			"vBulletin \u2265 v3.8.5",
		},
	},
	{
		regex: "^(\\$snefru\\$)?[a-f0-9]{64}$",
		names: []string{
			"Snefru-256",
		},
	},
	{
		regex: "^[a-f0-9]{64}(:.+)?$",
		names: []string{
			"SHA-256",
			"RIPEMD-256",
			"Haval-256",
			"GOST R 34.11-94",
			"GOST CryptoPro S-Box",
			"SHA3-256",
			"Skein-256",
			"Skein-512(256)",
			"Ventrilo",
			"sha256($pass.$salt)",
			"sha256($salt.$pass)",
			"sha256(unicode($pass).$salt)",
			"sha256($salt.unicode($pass))",
			"HMAC-SHA256 (key = $pass)",
			"HMAC-SHA256 (key = $salt)",
		},
	},
	{
		regex: "^[a-f0-9]{32}:[a-z0-9]{32}$",
		names: []string{
			"Joomla < v2.5.18",
		},
	},
	{
		regex: "^[a-f-0-9]{32}:[a-f-0-9]{32}$",
		names: []string{
			"SAM(LM_Hash:NT_Hash)",
		},
	},
	{
		regex: "^(\\$chap\\$0\\*)?[a-f0-9]{32}[\\*:][a-f0-9]{32}(:[0-9]{2})?$",
		names: []string{
			"MD5(Chap)",
			"iSCSI CHAP Authentication",
		},
	},
	{
		regex: "^\\$episerver\\$\\*0\\*[a-z0-9\\/=+]+\\*[a-z0-9\\/=+]{27,28}$",
		names: []string{
			"EPiServer 6.x < v4",
		},
	},
	{
		regex: "^{ssha256}[0-9]{2}\\$[a-z0-9$\\/.]{60}$",
		names: []string{
			"AIX(ssha256)",
		},
	},
	{
		regex: "^[a-f0-9]{80}$",
		names: []string{
			"RIPEMD-320",
		},
	},
	{
		regex: "^\\$episerver\\$\\*1\\*[a-z0-9\\/=+]+\\*[a-z0-9\\/=+]{42,43}$",
		names: []string{
			"EPiServer 6.x \u2265 v4",
		},
	},
	{
		regex: "^0x0100[a-f0-9]{88}$",
		names: []string{
			"MSSQL(2000)",
		},
	},
	{
		regex: "^[a-f0-9]{96}$",
		names: []string{
			"SHA-384",
			"SHA3-384",
			"Skein-512(384)",
			"Skein-1024(384)",
		},
	},
	{
		regex: "^{SSHA512}[a-z0-9\\/+]{96}$",
		names: []string{
			"SSHA-512(Base64)",
			"LDAP(SSHA-512)",
		},
	},
	{
		regex: "^{ssha512}[0-9]{2}\\$[a-z0-9\\/.]{16,48}\\$[a-z0-9\\/.]{86}$",
		names: []string{
			"AIX(ssha512)",
		},
	},
	{
		regex: "^[a-f0-9]{128}(:.+)?$",
		names: []string{
			"SHA-512",
			"Whirlpool",
			"Salsa10",
			"Salsa20",
			"SHA3-512",
			"Skein-512",
			"Skein-1024(512)",
			"sha512($pass.$salt)",
			"sha512($salt.$pass)",
			"sha512(unicode($pass).$salt)",
			"sha512($salt.unicode($pass))",
			"HMAC-SHA512 (key = $pass)",
			"HMAC-SHA512 (key = $salt)",
		},
	},
	{
		regex: "^[a-f0-9]{136}$",
		names: []string{
			"OSX v10.7",
		},
	},
	{
		regex: "^0x0200[a-f0-9]{136}$",
		names: []string{
			"MSSQL(2012)",
			"MSSQL(2014)",
		},
	},
	{
		regex: "^\\$ml\\$[0-9]+\\$[a-f0-9]{64}\\$[a-f0-9]{128}$",
		names: []string{
			"OSX v10.8",
			"OSX v10.9",
		},
	},
	{
		regex: "^[a-f0-9]{256}$",
		names: []string{
			"Skein-1024",
		},
	},
	{
		regex: "^grub\\.pbkdf2\\.sha512\\.[0-9]+\\.([a-f0-9]{128,1000}[a-f0-9]{0,1000}[a-f0-9]{0,48}\\.|[0-9]+\\.)?[a-f0-9]{128}$",
		names: []string{
			"GRUB 2",
		},
	},
	{
		regex: "^sha1\\$[a-z0-9]+\\$[a-f0-9]{40}$",
		names: []string{
			"Django(SHA-1)",
		},
	},
	{
		regex: "^[a-f0-9]{49}$",
		names: []string{
			"Citrix Netscaler",
		},
	},
	{
		regex: "^\\$S\\$[a-z0-9\\/.]{52}$",
		names: []string{
			"Drupal > v7.x",
		},
	},
	{
		regex: "^\\$5\\$(rounds=[0-9]+\\$)?[a-z0-9\\/.]{0,16}\\$[a-z0-9\\/.]{43}$",
		names: []string{
			"SHA-256 Crypt",
		},
	},
	{
		regex: "^0x[a-f0-9]{4}[a-f0-9]{16}[a-f0-9]{64}$",
		names: []string{
			"Sybase ASE",
		},
	},
	{
		regex: "^\\$6\\$(rounds=[0-9]+\\$)?[a-z0-9\\/.]{0,16}\\$[a-z0-9\\/.]{86}$",
		names: []string{
			"SHA-512 Crypt",
		},
	},
	{
		regex: "^\\$sha\\$[a-z0-9]{1,16}\\$([a-f0-9]{32}|[a-f0-9]{40}|[a-f0-9]{64}|[a-f0-9]{128}|[a-f0-9]{140})$",
		names: []string{
			"Minecraft(AuthMe Reloaded)",
		},
	},
	{
		regex: "^sha256\\$[a-z0-9]+\\$[a-f0-9]{64}$",
		names: []string{
			"Django(SHA-256)",
		},
	},
	{
		regex: "^sha384\\$[a-z0-9]+\\$[a-f0-9]{96}$",
		names: []string{
			"Django(SHA-384)",
		},
	},
	{
		regex: "^crypt1:[a-z0-9+=]{12}:[a-z0-9+=]{12}$",
		names: []string{
			"Clavister Secure Gateway",
		},
	},
	{
		regex: "^[a-f0-9]{112}$",
		names: []string{
			"Cisco VPN Client(PCF-File)",
		},
	},
	{
		regex: "^[a-f0-9]{1000}[a-f0-9]{329}$",
		names: []string{
			"Microsoft MSTSC(RDP-File)",
		},
	},
	{
		regex: "^[^\\\\\\/:*?\"<>|]{1,20}[:]{2,3}([^\\\\\\/:*?\"<>|]{1,20})?:[a-f0-9]{48}:[a-f0-9]{48}:[a-f0-9]{16}$",
		names: []string{
			"NetNTLMv1-VANILLA / NetNTLMv1+ESS",
		},
	},
	{
		regex: "^([^\\\\\\/:*?\"<>|]{1,20}\\\\)?[^\\\\\\/:*?\"<>|]{1,20}[:]{2,3}([^\\\\\\/:*?\"<>|]{1,20}:)?[^\\\\\\/:*?\"<>|]{1,20}:[a-f0-9]{32}:[a-f0-9]+$",
		names: []string{
			"NetNTLMv2",
		},
	},
	{
		regex: "^\\$(krb5pa|mskrb5)\\$([0-9]{2})?\\$.+\\$[a-f0-9]{1,}$",
		names: []string{
			"Kerberos 5 AS-REQ Pre-Auth",
		},
	},
	{
		regex: "^\\$scram\\$[0-9]+\\$[a-z0-9\\/.]{16}\\$sha-1=[a-z0-9\\/.]{27},sha-256=[a-z0-9\\/.]{43},sha-512=[a-z0-9\\/.]{86}$",
		names: []string{
			"SCRAM Hash",
		},
	},
	{
		regex: "^[a-f0-9]{40}:[a-f0-9]{0,32}$",
		names: []string{
			"Redmine Project Management Web App",
		},
	},
	{
		regex: "^(.+)?\\$[a-f0-9]{16}$",
		names: []string{
			"SAP CODVN B (BCODE)",
		},
	},
	{
		regex: "^(.+)?\\$[a-f0-9]{40}$",
		names: []string{
			"SAP CODVN F/G (PASSCODE)",
		},
	},
	{
		regex: "^(.+\\$)?[a-z0-9\\/.+]{30}(:.+)?$",
		names: []string{
			"Juniper Netscreen/SSG(ScreenOS)",
		},
	},
	{
		regex: "^0x[a-f0-9]{60}\\s0x[a-f0-9]{40}$",
		names: []string{
			"EPi",
		},
	},
	{
		regex: "^[a-f0-9]{40}:[^*]{1,25}$",
		names: []string{
			"SMF \u2265 v1.1",
		},
	},
	{
		regex: "^(\\$wbb3\\$\\*1\\*)?[a-f0-9]{40}[:*][a-f0-9]{40}$",
		names: []string{
			"Woltlab Burning Board 3.x",
		},
	},
	{
		regex: "^[a-f0-9]{130}(:[a-f0-9]{40})?$",
		names: []string{
			"IPMI2 RAKP HMAC-SHA1",
		},
	},
	{
		regex: "^[a-f0-9]{32}:[0-9]+:[a-z0-9_.+-]+@[a-z0-9-]+\\.[a-z0-9-.]+$",
		names: []string{
			"Lastpass",
		},
	},
	{
		regex: "^[a-z0-9\\/.]{16}([:$].{1,})?$",
		names: []string{
			"Cisco-ASA(MD5)",
		},
	},
	{
		regex: "^\\$vnc\\$\\*[a-f0-9]{32}\\*[a-f0-9]{32}$",
		names: []string{
			"VNC",
		},
	},
	{
		regex: "^[a-z0-9]{32}(:([a-z0-9-]+\\.)?[a-z0-9-.]+\\.[a-z]{2,7}:.+:[0-9]+)?$",
		names: []string{
			"DNSSEC(NSEC3)",
		},
	},
	{
		regex: "^(user-.+:)?\\$racf\\$\\*.+\\*[a-f0-9]{16}$",
		names: []string{
			"RACF",
		},
	},
	{
		regex: "^\\$3\\$\\$[a-f0-9]{32}$",
		names: []string{
			"NTHash(FreeBSD Variant)",
		},
	},
	{
		regex: "^\\$sha1\\$[0-9]+\\$[a-z0-9\\/.]{0,64}\\$[a-z0-9\\/.]{28}$",
		names: []string{
			"SHA-1 Crypt",
		},
	},
	{
		regex: "^[a-f0-9]{70}$",
		names: []string{
			"hMailServer",
		},
	},
	{
		regex: "^[:\\$][AB][:\\$]([a-f0-9]{1,8}[:\\$])?[a-f0-9]{32}$",
		names: []string{
			"MediaWiki",
		},
	},
	{
		regex: "^[a-f0-9]{140}$",
		names: []string{
			"Minecraft(xAuth)",
		},
	},
	{
		regex: "^\\$pbkdf2(-sha1)?\\$[0-9]+\\$[a-z0-9\\/.]+\\$[a-z0-9\\/.]{27}$",
		names: []string{
			"PBKDF2-SHA1(Generic)",
		},
	},
	{
		regex: "^\\$pbkdf2-sha256\\$[0-9]+\\$[a-z0-9\\/.]+\\$[a-z0-9\\/.]{43}$",
		names: []string{
			"PBKDF2-SHA256(Generic)",
		},
	},
	{
		regex: "^\\$pbkdf2-sha512\\$[0-9]+\\$[a-z0-9\\/.]+\\$[a-z0-9\\/.]{86}$",
		names: []string{
			"PBKDF2-SHA512(Generic)",
		},
	},
	{
		regex: "^\\$p5k2\\$[0-9]+\\$[a-z0-9\\/+=-]+\\$[a-z0-9\\/+-]{27}=$",
		names: []string{
			"PBKDF2(Cryptacular)",
		},
	},
	{
		regex: "^\\$p5k2\\$[0-9]+\\$[a-z0-9\\/.]+\\$[a-z0-9\\/.]{32}$",
		names: []string{
			"PBKDF2(Dwayne Litzenberger)",
		},
	},
	{
		regex: "^{FSHP[0123]\\|[0-9]+\\|[0-9]+}[a-z0-9\\/+=]+$",
		names: []string{
			"Fairly Secure Hashed Password",
		},
	},
	{
		regex: "^\\$PHPS\\$.+\\$[a-f0-9]{32}$",
		names: []string{
			"PHPS",
		},
	},
	{
		regex: "^[0-9]{4}:[a-f0-9]{16}:[a-f0-9]{1000}[a-f0-9]{1000}[a-f0-9]{80}$",
		names: []string{
			"1Password(Agile Keychain)",
		},
	},
	{
		regex: "^[a-f0-9]{64}:[a-f0-9]{32}:[0-9]{5}:[a-f0-9]{608}$",
		names: []string{
			"1Password(Cloud Keychain)",
		},
	},
	{
		regex: "^[a-f0-9]{256}:[a-f0-9]{256}:[a-f0-9]{16}:[a-f0-9]{16}:[a-f0-9]{320}:[a-f0-9]{16}:[a-f0-9]{40}:[a-f0-9]{40}:[a-f0-9]{32}$",
		names: []string{
			"IKE-PSK MD5",
		},
	},
	{
		regex: "^[a-f0-9]{256}:[a-f0-9]{256}:[a-f0-9]{16}:[a-f0-9]{16}:[a-f0-9]{320}:[a-f0-9]{16}:[a-f0-9]{40}:[a-f0-9]{40}:[a-f0-9]{40}$",
		names: []string{
			"IKE-PSK SHA1",
		},
	},
	{
		regex: "^[a-z0-9\\/+]{27}=$",
		names: []string{
			"PeopleSoft",
		},
	},
	{
		regex: "^crypt\\$[a-f0-9]{5}\\$[a-z0-9\\/.]{13}$",
		names: []string{
			"Django(DES Crypt Wrapper)",
		},
	},
	{
		regex: "^(\\$django\\$\\*1\\*)?pbkdf2_sha256\\$[0-9]+\\$[a-z0-9]+\\$[a-z0-9\\/+=]{44}$",
		names: []string{
			"Django(PBKDF2-HMAC-SHA256)",
		},
	},
	{
		regex: "^pbkdf2_sha1\\$[0-9]+\\$[a-z0-9]+\\$[a-z0-9\\/+=]{28}$",
		names: []string{
			"Django(PBKDF2-HMAC-SHA1)",
		},
	},
	{
		regex: "^bcrypt(\\$2[axy]|\\$2)\\$[0-9]{2}\\$[a-z0-9\\/.]{53}$",
		names: []string{
			"Django(bcrypt)",
		},
	},
	{
		regex: "^md5\\$[a-f0-9]+\\$[a-f0-9]{32}$",
		names: []string{
			"Django(MD5)",
		},
	},
	{
		regex: "^\\{PKCS5S2\\}[a-z0-9\\/+]{64}$",
		names: []string{
			"PBKDF2(Atlassian)",
		},
	},
	{
		regex: "^md5[a-f0-9]{32}$",
		names: []string{
			"PostgreSQL MD5",
		},
	},
	{
		regex: "^\\([a-z0-9\\/+]{49}\\)$",
		names: []string{
			"Lotus Notes/Domino 8",
		},
	},
	{
		regex: "^SCRYPT:[0-9]{1,}:[0-9]{1}:[0-9]{1}:[a-z0-9:\\/+=]{1,}$",
		names: []string{
			"scrypt",
		},
	},
	{
		regex: "^\\$8\\$[a-z0-9\\/.]{14}\\$[a-z0-9\\/.]{43}$",
		names: []string{
			"Cisco Type 8",
		},
	},
	{
		regex: "^\\$9\\$[a-z0-9\\/.]{14}\\$[a-z0-9\\/.]{43}$",
		names: []string{
			"Cisco Type 9",
		},
	},
	{
		regex: "^\\$office\\$\\*2007\\*[0-9]{2}\\*[0-9]{3}\\*[0-9]{2}\\*[a-z0-9]{32}\\*[a-z0-9]{32}\\*[a-z0-9]{40}$",
		names: []string{
			"Microsoft Office 2007",
		},
	},
	{
		regex: "^\\$office\\$\\*2010\\*[0-9]{6}\\*[0-9]{3}\\*[0-9]{2}\\*[a-z0-9]{32}\\*[a-z0-9]{32}\\*[a-z0-9]{64}$",
		names: []string{
			"Microsoft Office 2010",
		},
	},
	{
		regex: "^\\$office\\$\\*2013\\*[0-9]{6}\\*[0-9]{3}\\*[0-9]{2}\\*[a-z0-9]{32}\\*[a-z0-9]{32}\\*[a-z0-9]{64}$",
		names: []string{
			"Microsoft Office 2013",
		},
	},
	{
		regex: "^\\$fde\\$[0-9]{2}\\$[a-f0-9]{32}\\$[0-9]{2}\\$[a-f0-9]{32}\\$[a-f0-9]{1000}[a-f0-9]{1000}[a-f0-9]{1000}[a-f0-9]{72}$",
		names: []string{
			"Android FDE \u2264 4.3",
		},
	},
	{
		regex: "^\\$oldoffice\\$[01]\\*[a-f0-9]{32}\\*[a-f0-9]{32}\\*[a-f0-9]{32}$",
		names: []string{
			"Microsoft Office \u2264 2003 (MD5+RC4)",
			"Microsoft Office \u2264 2003 (MD5+RC4) collider-mode #1",
			"Microsoft Office \u2264 2003 (MD5+RC4) collider-mode #2",
		},
	},
	{
		regex: "^\\$oldoffice\\$[34]\\*[a-f0-9]{32}\\*[a-f0-9]{32}\\*[a-f0-9]{40}$",
		names: []string{
			"Microsoft Office \u2264 2003 (SHA1+RC4)",
			"Microsoft Office \u2264 2003 (SHA1+RC4) collider-mode #1",
			"Microsoft Office \u2264 2003 (SHA1+RC4) collider-mode #2",
		},
	},
	{
		regex: "^(\\$radmin2\\$)?[a-f0-9]{32}$",
		names: []string{
			"RAdmin v2.x",
		},
	},
	{
		regex: "^{x-issha,\\s[0-9]{4}}[a-z0-9\\/+=]+$",
		names: []string{
			"SAP CODVN H (PWDSALTEDHASH) iSSHA-1",
		},
	},
	{
		regex: "^\\$cram_md5\\$[a-z0-9\\/+=-]+\\$[a-z0-9\\/+=-]{52}$",
		names: []string{
			"CRAM-MD5",
		},
	},
	{
		regex: "^[a-f0-9]{16}:2:4:[a-f0-9]{32}$",
		names: []string{
			"SipHash",
		},
	},
	{
		regex: "^[a-f0-9]{4,}$",
		names: []string{
			"Cisco Type 7",
		},
	},
	{
		regex: "^[a-z0-9\\/.]{13,}$",
		names: []string{
			"BigCrypt",
		},
	},
	{
		regex: "^(\\$cisco4\\$)?[a-z0-9\\/.]{43}$",
		names: []string{
			"Cisco Type 4",
		},
	},
	{
		regex: "^bcrypt_sha256\\$\\$(2[axy]|2)\\$[0-9]+\\$[a-z0-9\\/.]{53}$",
		names: []string{
			"Django(bcrypt-SHA256)",
		},
	},
	{
		regex: "^\\$postgres\\$.[^\\*]+[*:][a-f0-9]{1,32}[*:][a-f0-9]{32}$",
		names: []string{
			"PostgreSQL Challenge-Response Authentication (MD5)",
		},
	},
	{
		regex: "^\\$siemens-s7\\$[0-9]{1}\\$[a-f0-9]{40}\\$[a-f0-9]{40}$",
		names: []string{
			"Siemens-S7",
		},
	},
	{
		regex: "^(\\$pst\\$)?[a-f0-9]{8}$",
		names: []string{
			"Microsoft Outlook PST",
		},
	},
	{
		regex: "^sha256[:$][0-9]+[:$][a-z0-9\\/+]+[:$][a-z0-9\\/+]{32,128}$",
		names: []string{
			"PBKDF2-HMAC-SHA256(PHP)",
		},
	},
	{
		regex: "^(\\$dahua\\$)?[a-z0-9]{8}$",
		names: []string{
			"Dahua",
		},
	},
	{
		regex: "^\\$mysqlna\\$[a-f0-9]{40}[:*][a-f0-9]{40}$",
		names: []string{
			"MySQL Challenge-Response Authentication (SHA1)",
		},
	},
	{
		regex: "^\\$pdf\\$[24]\\*[34]\\*128\\*[0-9-]{1,5}\\*1\\*(16|32)\\*[a-f0-9]{32,64}\\*32\\*[a-f0-9]{64}\\*(8|16|32)\\*[a-f0-9]{16,64}$",
		names: []string{
			"PDF 1.4 - 1.6 (Acrobat 5 - 8)",
		},
	},
}

func Lookup(s string) <-chan string {
	ch := make(chan string)

	go func() {
		for _, e := range db {
			re := regexp.MustCompile(e.regex)

			if re.MatchString(s) {
				for _, s := range e.names {
					ch <- s
				}
			}
		}

		close(ch)
	}()

	return ch
}
