// https://github.com/s0md3v/Bolt/blob/master/db/hashes.json
package text

type Pattern struct {
	Regex string
	Names []string
}

var Patterns = []Pattern{
	{
		Regex: "^[a-f0-9]{4}$",
		Names: []string{
			"CRC-16",
			"CRC-16-CCITT",
			"FCS-16",
		},
	},
	{
		Regex: "^[a-f0-9]{8}$",
		Names: []string{
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
		Regex: "^[a-f0-9]{6}$",
		Names: []string{
			"CRC-24",
		},
	},
	{
		Regex: "^(\\$crc32\\$[a-f0-9]{8}.)?[a-f0-9]{8}$",
		Names: []string{
			"CRC-32",
		},
	},
	{
		Regex: "^\\+[a-z0-9\\/.]{12}$",
		Names: []string{
			"Eggdrop IRC Bot",
		},
	},
	{
		Regex: "^[a-z0-9\\/.]{13}$",
		Names: []string{
			"DES(Unix)",
			"Traditional DES",
			"DEScrypt",
		},
	},
	{
		Regex: "^[a-f0-9]{16}$",
		Names: []string{
			"MySQL323",
			"DES(Oracle)",
			"Half MD5",
			"Oracle 7-10g",
			"FNV-164",
			"CRC-64",
		},
	},
	{
		Regex: "^[a-z0-9\\/.]{16}$",
		Names: []string{
			"Cisco-PIX(MD5)",
		},
	},
	{
		Regex: "^\\([a-z0-9\\/+]{20}\\)$",
		Names: []string{
			"Lotus Notes/Domino 6",
		},
	},
	{
		Regex: "^_[a-z0-9\\/.]{19}$",
		Names: []string{
			"BSDi Crypt",
		},
	},
	{
		Regex: "^[a-f0-9]{24}$",
		Names: []string{
			"CRC-96(ZIP)",
		},
	},
	{
		Regex: "^[a-z0-9\\/.]{24}$",
		Names: []string{
			"Crypt16",
		},
	},
	{
		Regex: "^(\\$md2\\$)?[a-f0-9]{32}$",
		Names: []string{
			"MD2",
		},
	},
	{
		Regex: "^[a-f0-9]{32}(:.+)?$",
		Names: []string{
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
		Regex: "^(\\$snefru\\$)?[a-f0-9]{32}$",
		Names: []string{
			"Snefru-128",
		},
	},
	{
		Regex: "^(\\$NT\\$)?[a-f0-9]{32}$",
		Names: []string{
			"NTLM",
		},
	},
	{
		Regex: "^([^\\\\\\/:*?\"<>|]{1,20}:)?[a-f0-9]{32}(:[^\\\\\\/:*?\"<>|]{1,20})?$",
		Names: []string{
			"Domain Cached Credentials",
		},
	},
	{
		Regex: "^([^\\\\\\/:*?\"<>|]{1,20}:)?(\\$DCC2\\$10240#[^\\\\\\/:*?\"<>|]{1,20}#)?[a-f0-9]{32}$",
		Names: []string{
			"Domain Cached Credentials 2",
		},
	},
	{
		Regex: "^{SHA}[a-z0-9\\/+]{27}=$",
		Names: []string{
			"SHA-1(Base64)",
			"Netscape LDAP SHA",
		},
	},
	{
		Regex: "^\\$1\\$[a-z0-9\\/.]{0,8}\\$[a-z0-9\\/.]{22}(:.*)?$",
		Names: []string{
			"MD5 Crypt",
			"Cisco-IOS(MD5)",
			"FreeBSD MD5",
		},
	},
	{
		Regex: "^0x[a-f0-9]{32}$",
		Names: []string{
			"Lineage II C4",
		},
	},
	{
		Regex: "^\\$H\\$[a-z0-9\\/.]{31}$",
		Names: []string{
			"phpBB v3.x",
			"Wordpress v2.6.0/2.6.1",
			"PHPass' Portable Hash",
		},
	},
	{
		Regex: "^\\$P\\$[a-z0-9\\/.]{31}$",
		Names: []string{
			"Wordpress \u2265 v2.6.2",
			"Joomla \u2265 v2.5.18",
			"PHPass' Portable Hash",
		},
	},
	{
		Regex: "^[a-f0-9]{32}:[a-z0-9]{2}$",
		Names: []string{
			"osCommerce",
			"xt:Commerce",
		},
	},
	{
		Regex: "^\\$apr1\\$[a-z0-9\\/.]{0,8}\\$[a-z0-9\\/.]{22}$",
		Names: []string{
			"MD5(APR)",
			"Apache MD5",
			"md5apr1",
		},
	},
	{
		Regex: "^{smd5}[a-z0-9$\\/.]{31}$",
		Names: []string{
			"AIX(smd5)",
		},
	},
	{
		Regex: "^[a-f0-9]{32}:[a-f0-9]{32}$",
		Names: []string{
			"WebEdition CMS",
		},
	},
	{
		Regex: "^[a-f0-9]{32}:.{5}$",
		Names: []string{
			"IP.Board \u2265 v2+",
		},
	},
	{
		Regex: "^[a-f0-9]{32}:.{8}$",
		Names: []string{
			"MyBB \u2265 v1.2+",
		},
	},
	{
		Regex: "^[a-z0-9]{34}$",
		Names: []string{
			"CryptoCurrency(Address)",
		},
	},
	{
		Regex: "^[a-f0-9]{40}(:.+)?$",
		Names: []string{
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
		Regex: "^\\*[a-f0-9]{40}$",
		Names: []string{
			"MySQL5.x",
			"MySQL4.1",
		},
	},
	{
		Regex: "^[a-z0-9]{43}$",
		Names: []string{
			"Cisco-IOS(SHA-256)",
		},
	},
	{
		Regex: "^{SSHA}[a-z0-9\\/+]{38}==$",
		Names: []string{
			"SSHA-1(Base64)",
			"Netscape LDAP SSHA",
			"nsldaps",
		},
	},
	{
		Regex: "^[a-z0-9=]{47}$",
		Names: []string{
			"Fortigate(FortiOS)",
		},
	},
	{
		Regex: "^[a-f0-9]{48}$",
		Names: []string{
			"Haval-192",
			"Tiger-192",
			"SHA-1(Oracle)",
			"OSX v10.4",
			"OSX v10.5",
			"OSX v10.6",
		},
	},
	{
		Regex: "^[a-f0-9]{51}$",
		Names: []string{
			"Palshop CMS",
		},
	},
	{
		Regex: "^[a-z0-9]{51}$",
		Names: []string{
			"CryptoCurrency(PrivateKey)",
		},
	},
	{
		Regex: "^{ssha1}[0-9]{2}\\$[a-z0-9$\\/.]{44}$",
		Names: []string{
			"AIX(ssha1)",
		},
	},
	{
		Regex: "^0x0100[a-f0-9]{48}$",
		Names: []string{
			"MSSQL(2005)",
			"MSSQL(2008)",
		},
	},
	{
		Regex: "^(\\$md5,rounds=[0-9]+\\$|\\$md5\\$rounds=[0-9]+\\$|\\$md5\\$)[a-z0-9\\/.]{0,16}(\\$|\\$\\$)[a-z0-9\\/.]{22}$",
		Names: []string{
			"Sun MD5 Crypt",
		},
	},
	{
		Regex: "^[a-f0-9]{56}$",
		Names: []string{
			"SHA-224",
			"Haval-224",
			"SHA3-224",
			"Skein-256(224)",
			"Skein-512(224)",
		},
	},
	{
		Regex: "^(\\$2[axy]|\\$2)\\$[0-9]{2}\\$[a-z0-9\\/.]{53}$",
		Names: []string{
			"Blowfish(OpenBSD)",
			"Woltlab Burning Board 4.x",
			"bcrypt",
		},
	},
	{
		Regex: "^[a-f0-9]{40}:[a-f0-9]{16}$",
		Names: []string{
			"Android PIN",
		},
	},
	{
		Regex: "^(S:)?[a-f0-9]{40}(:)?[a-f0-9]{20}$",
		Names: []string{
			"Oracle 11g/12c",
		},
	},
	{
		Regex: "^\\$bcrypt-sha256\\$(2[axy]|2)\\,[0-9]+\\$[a-z0-9\\/.]{22}\\$[a-z0-9\\/.]{31}$",
		Names: []string{
			"bcrypt(SHA-256)",
		},
	},
	{
		Regex: "^[a-f0-9]{32}:.{3}$",
		Names: []string{
			"vBulletin < v3.8.5",
		},
	},
	{
		Regex: "^[a-f0-9]{32}:.{30}$",
		Names: []string{
			"vBulletin \u2265 v3.8.5",
		},
	},
	{
		Regex: "^(\\$snefru\\$)?[a-f0-9]{64}$",
		Names: []string{
			"Snefru-256",
		},
	},
	{
		Regex: "^[a-f0-9]{64}(:.+)?$",
		Names: []string{
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
		Regex: "^[a-f0-9]{32}:[a-z0-9]{32}$",
		Names: []string{
			"Joomla < v2.5.18",
		},
	},
	{
		Regex: "^[a-f-0-9]{32}:[a-f-0-9]{32}$",
		Names: []string{
			"SAM(LM_Hash:NT_Hash)",
		},
	},
	{
		Regex: "^(\\$chap\\$0\\*)?[a-f0-9]{32}[\\*:][a-f0-9]{32}(:[0-9]{2})?$",
		Names: []string{
			"MD5(Chap)",
			"iSCSI CHAP Authentication",
		},
	},
	{
		Regex: "^\\$episerver\\$\\*0\\*[a-z0-9\\/=+]+\\*[a-z0-9\\/=+]{27,28}$",
		Names: []string{
			"EPiServer 6.x < v4",
		},
	},
	{
		Regex: "^{ssha256}[0-9]{2}\\$[a-z0-9$\\/.]{60}$",
		Names: []string{
			"AIX(ssha256)",
		},
	},
	{
		Regex: "^[a-f0-9]{80}$",
		Names: []string{
			"RIPEMD-320",
		},
	},
	{
		Regex: "^\\$episerver\\$\\*1\\*[a-z0-9\\/=+]+\\*[a-z0-9\\/=+]{42,43}$",
		Names: []string{
			"EPiServer 6.x \u2265 v4",
		},
	},
	{
		Regex: "^0x0100[a-f0-9]{88}$",
		Names: []string{
			"MSSQL(2000)",
		},
	},
	{
		Regex: "^[a-f0-9]{96}$",
		Names: []string{
			"SHA-384",
			"SHA3-384",
			"Skein-512(384)",
			"Skein-1024(384)",
		},
	},
	{
		Regex: "^{SSHA512}[a-z0-9\\/+]{96}$",
		Names: []string{
			"SSHA-512(Base64)",
			"LDAP(SSHA-512)",
		},
	},
	{
		Regex: "^{ssha512}[0-9]{2}\\$[a-z0-9\\/.]{16,48}\\$[a-z0-9\\/.]{86}$",
		Names: []string{
			"AIX(ssha512)",
		},
	},
	{
		Regex: "^[a-f0-9]{128}(:.+)?$",
		Names: []string{
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
		Regex: "^[a-f0-9]{136}$",
		Names: []string{
			"OSX v10.7",
		},
	},
	{
		Regex: "^0x0200[a-f0-9]{136}$",
		Names: []string{
			"MSSQL(2012)",
			"MSSQL(2014)",
		},
	},
	{
		Regex: "^\\$ml\\$[0-9]+\\$[a-f0-9]{64}\\$[a-f0-9]{128}$",
		Names: []string{
			"OSX v10.8",
			"OSX v10.9",
		},
	},
	{
		Regex: "^[a-f0-9]{256}$",
		Names: []string{
			"Skein-1024",
		},
	},
	{
		Regex: "^grub\\.pbkdf2\\.sha512\\.[0-9]+\\.([a-f0-9]{128,1000}[a-f0-9]{0,1000}[a-f0-9]{0,48}\\.|[0-9]+\\.)?[a-f0-9]{128}$",
		Names: []string{
			"GRUB 2",
		},
	},
	{
		Regex: "^sha1\\$[a-z0-9]+\\$[a-f0-9]{40}$",
		Names: []string{
			"Django(SHA-1)",
		},
	},
	{
		Regex: "^[a-f0-9]{49}$",
		Names: []string{
			"Citrix Netscaler",
		},
	},
	{
		Regex: "^\\$S\\$[a-z0-9\\/.]{52}$",
		Names: []string{
			"Drupal > v7.x",
		},
	},
	{
		Regex: "^\\$5\\$(rounds=[0-9]+\\$)?[a-z0-9\\/.]{0,16}\\$[a-z0-9\\/.]{43}$",
		Names: []string{
			"SHA-256 Crypt",
		},
	},
	{
		Regex: "^0x[a-f0-9]{4}[a-f0-9]{16}[a-f0-9]{64}$",
		Names: []string{
			"Sybase ASE",
		},
	},
	{
		Regex: "^\\$6\\$(rounds=[0-9]+\\$)?[a-z0-9\\/.]{0,16}\\$[a-z0-9\\/.]{86}$",
		Names: []string{
			"SHA-512 Crypt",
		},
	},
	{
		Regex: "^\\$sha\\$[a-z0-9]{1,16}\\$([a-f0-9]{32}|[a-f0-9]{40}|[a-f0-9]{64}|[a-f0-9]{128}|[a-f0-9]{140})$",
		Names: []string{
			"Minecraft(AuthMe Reloaded)",
		},
	},
	{
		Regex: "^sha256\\$[a-z0-9]+\\$[a-f0-9]{64}$",
		Names: []string{
			"Django(SHA-256)",
		},
	},
	{
		Regex: "^sha384\\$[a-z0-9]+\\$[a-f0-9]{96}$",
		Names: []string{
			"Django(SHA-384)",
		},
	},
	{
		Regex: "^crypt1:[a-z0-9+=]{12}:[a-z0-9+=]{12}$",
		Names: []string{
			"Clavister Secure Gateway",
		},
	},
	{
		Regex: "^[a-f0-9]{112}$",
		Names: []string{
			"Cisco VPN Client(PCF-File)",
		},
	},
	{
		Regex: "^[a-f0-9]{1000}[a-f0-9]{329}$",
		Names: []string{
			"Microsoft MSTSC(RDP-File)",
		},
	},
	{
		Regex: "^[^\\\\\\/:*?\"<>|]{1,20}[:]{2,3}([^\\\\\\/:*?\"<>|]{1,20})?:[a-f0-9]{48}:[a-f0-9]{48}:[a-f0-9]{16}$",
		Names: []string{
			"NetNTLMv1-VANILLA / NetNTLMv1+ESS",
		},
	},
	{
		Regex: "^([^\\\\\\/:*?\"<>|]{1,20}\\\\)?[^\\\\\\/:*?\"<>|]{1,20}[:]{2,3}([^\\\\\\/:*?\"<>|]{1,20}:)?[^\\\\\\/:*?\"<>|]{1,20}:[a-f0-9]{32}:[a-f0-9]+$",
		Names: []string{
			"NetNTLMv2",
		},
	},
	{
		Regex: "^\\$(krb5pa|mskrb5)\\$([0-9]{2})?\\$.+\\$[a-f0-9]{1,}$",
		Names: []string{
			"Kerberos 5 AS-REQ Pre-Auth",
		},
	},
	{
		Regex: "^\\$scram\\$[0-9]+\\$[a-z0-9\\/.]{16}\\$sha-1=[a-z0-9\\/.]{27},sha-256=[a-z0-9\\/.]{43},sha-512=[a-z0-9\\/.]{86}$",
		Names: []string{
			"SCRAM Hash",
		},
	},
	{
		Regex: "^[a-f0-9]{40}:[a-f0-9]{0,32}$",
		Names: []string{
			"Redmine Project Management Web App",
		},
	},
	{
		Regex: "^(.+)?\\$[a-f0-9]{16}$",
		Names: []string{
			"SAP CODVN B (BCODE)",
		},
	},
	{
		Regex: "^(.+)?\\$[a-f0-9]{40}$",
		Names: []string{
			"SAP CODVN F/G (PASSCODE)",
		},
	},
	{
		Regex: "^(.+\\$)?[a-z0-9\\/.+]{30}(:.+)?$",
		Names: []string{
			"Juniper Netscreen/SSG(ScreenOS)",
		},
	},
	{
		Regex: "^0x[a-f0-9]{60}\\s0x[a-f0-9]{40}$",
		Names: []string{
			"EPi",
		},
	},
	{
		Regex: "^[a-f0-9]{40}:[^*]{1,25}$",
		Names: []string{
			"SMF \u2265 v1.1",
		},
	},
	{
		Regex: "^(\\$wbb3\\$\\*1\\*)?[a-f0-9]{40}[:*][a-f0-9]{40}$",
		Names: []string{
			"Woltlab Burning Board 3.x",
		},
	},
	{
		Regex: "^[a-f0-9]{130}(:[a-f0-9]{40})?$",
		Names: []string{
			"IPMI2 RAKP HMAC-SHA1",
		},
	},
	{
		Regex: "^[a-f0-9]{32}:[0-9]+:[a-z0-9_.+-]+@[a-z0-9-]+\\.[a-z0-9-.]+$",
		Names: []string{
			"Lastpass",
		},
	},
	{
		Regex: "^[a-z0-9\\/.]{16}([:$].{1,})?$",
		Names: []string{
			"Cisco-ASA(MD5)",
		},
	},
	{
		Regex: "^\\$vnc\\$\\*[a-f0-9]{32}\\*[a-f0-9]{32}$",
		Names: []string{
			"VNC",
		},
	},
	{
		Regex: "^[a-z0-9]{32}(:([a-z0-9-]+\\.)?[a-z0-9-.]+\\.[a-z]{2,7}:.+:[0-9]+)?$",
		Names: []string{
			"DNSSEC(NSEC3)",
		},
	},
	{
		Regex: "^(user-.+:)?\\$racf\\$\\*.+\\*[a-f0-9]{16}$",
		Names: []string{
			"RACF",
		},
	},
	{
		Regex: "^\\$3\\$\\$[a-f0-9]{32}$",
		Names: []string{
			"NTHash(FreeBSD Variant)",
		},
	},
	{
		Regex: "^\\$sha1\\$[0-9]+\\$[a-z0-9\\/.]{0,64}\\$[a-z0-9\\/.]{28}$",
		Names: []string{
			"SHA-1 Crypt",
		},
	},
	{
		Regex: "^[a-f0-9]{70}$",
		Names: []string{
			"hMailServer",
		},
	},
	{
		Regex: "^[:\\$][AB][:\\$]([a-f0-9]{1,8}[:\\$])?[a-f0-9]{32}$",
		Names: []string{
			"MediaWiki",
		},
	},
	{
		Regex: "^[a-f0-9]{140}$",
		Names: []string{
			"Minecraft(xAuth)",
		},
	},
	{
		Regex: "^\\$pbkdf2(-sha1)?\\$[0-9]+\\$[a-z0-9\\/.]+\\$[a-z0-9\\/.]{27}$",
		Names: []string{
			"PBKDF2-SHA1(Generic)",
		},
	},
	{
		Regex: "^\\$pbkdf2-sha256\\$[0-9]+\\$[a-z0-9\\/.]+\\$[a-z0-9\\/.]{43}$",
		Names: []string{
			"PBKDF2-SHA256(Generic)",
		},
	},
	{
		Regex: "^\\$pbkdf2-sha512\\$[0-9]+\\$[a-z0-9\\/.]+\\$[a-z0-9\\/.]{86}$",
		Names: []string{
			"PBKDF2-SHA512(Generic)",
		},
	},
	{
		Regex: "^\\$p5k2\\$[0-9]+\\$[a-z0-9\\/+=-]+\\$[a-z0-9\\/+-]{27}=$",
		Names: []string{
			"PBKDF2(Cryptacular)",
		},
	},
	{
		Regex: "^\\$p5k2\\$[0-9]+\\$[a-z0-9\\/.]+\\$[a-z0-9\\/.]{32}$",
		Names: []string{
			"PBKDF2(Dwayne Litzenberger)",
		},
	},
	{
		Regex: "^{FSHP[0123]\\|[0-9]+\\|[0-9]+}[a-z0-9\\/+=]+$",
		Names: []string{
			"Fairly Secure Hashed Password",
		},
	},
	{
		Regex: "^\\$PHPS\\$.+\\$[a-f0-9]{32}$",
		Names: []string{
			"PHPS",
		},
	},
	{
		Regex: "^[0-9]{4}:[a-f0-9]{16}:[a-f0-9]{1000}[a-f0-9]{1000}[a-f0-9]{80}$",
		Names: []string{
			"1Password(Agile Keychain)",
		},
	},
	{
		Regex: "^[a-f0-9]{64}:[a-f0-9]{32}:[0-9]{5}:[a-f0-9]{608}$",
		Names: []string{
			"1Password(Cloud Keychain)",
		},
	},
	{
		Regex: "^[a-f0-9]{256}:[a-f0-9]{256}:[a-f0-9]{16}:[a-f0-9]{16}:[a-f0-9]{320}:[a-f0-9]{16}:[a-f0-9]{40}:[a-f0-9]{40}:[a-f0-9]{32}$",
		Names: []string{
			"IKE-PSK MD5",
		},
	},
	{
		Regex: "^[a-f0-9]{256}:[a-f0-9]{256}:[a-f0-9]{16}:[a-f0-9]{16}:[a-f0-9]{320}:[a-f0-9]{16}:[a-f0-9]{40}:[a-f0-9]{40}:[a-f0-9]{40}$",
		Names: []string{
			"IKE-PSK SHA1",
		},
	},
	{
		Regex: "^[a-z0-9\\/+]{27}=$",
		Names: []string{
			"PeopleSoft",
		},
	},
	{
		Regex: "^crypt\\$[a-f0-9]{5}\\$[a-z0-9\\/.]{13}$",
		Names: []string{
			"Django(DES Crypt Wrapper)",
		},
	},
	{
		Regex: "^(\\$django\\$\\*1\\*)?pbkdf2_sha256\\$[0-9]+\\$[a-z0-9]+\\$[a-z0-9\\/+=]{44}$",
		Names: []string{
			"Django(PBKDF2-HMAC-SHA256)",
		},
	},
	{
		Regex: "^pbkdf2_sha1\\$[0-9]+\\$[a-z0-9]+\\$[a-z0-9\\/+=]{28}$",
		Names: []string{
			"Django(PBKDF2-HMAC-SHA1)",
		},
	},
	{
		Regex: "^bcrypt(\\$2[axy]|\\$2)\\$[0-9]{2}\\$[a-z0-9\\/.]{53}$",
		Names: []string{
			"Django(bcrypt)",
		},
	},
	{
		Regex: "^md5\\$[a-f0-9]+\\$[a-f0-9]{32}$",
		Names: []string{
			"Django(MD5)",
		},
	},
	{
		Regex: "^\\{PKCS5S2\\}[a-z0-9\\/+]{64}$",
		Names: []string{
			"PBKDF2(Atlassian)",
		},
	},
	{
		Regex: "^md5[a-f0-9]{32}$",
		Names: []string{
			"PostgreSQL MD5",
		},
	},
	{
		Regex: "^\\([a-z0-9\\/+]{49}\\)$",
		Names: []string{
			"Lotus Notes/Domino 8",
		},
	},
	{
		Regex: "^SCRYPT:[0-9]{1,}:[0-9]{1}:[0-9]{1}:[a-z0-9:\\/+=]{1,}$",
		Names: []string{
			"scrypt",
		},
	},
	{
		Regex: "^\\$8\\$[a-z0-9\\/.]{14}\\$[a-z0-9\\/.]{43}$",
		Names: []string{
			"Cisco Type 8",
		},
	},
	{
		Regex: "^\\$9\\$[a-z0-9\\/.]{14}\\$[a-z0-9\\/.]{43}$",
		Names: []string{
			"Cisco Type 9",
		},
	},
	{
		Regex: "^\\$office\\$\\*2007\\*[0-9]{2}\\*[0-9]{3}\\*[0-9]{2}\\*[a-z0-9]{32}\\*[a-z0-9]{32}\\*[a-z0-9]{40}$",
		Names: []string{
			"Microsoft Office 2007",
		},
	},
	{
		Regex: "^\\$office\\$\\*2010\\*[0-9]{6}\\*[0-9]{3}\\*[0-9]{2}\\*[a-z0-9]{32}\\*[a-z0-9]{32}\\*[a-z0-9]{64}$",
		Names: []string{
			"Microsoft Office 2010",
		},
	},
	{
		Regex: "^\\$office\\$\\*2013\\*[0-9]{6}\\*[0-9]{3}\\*[0-9]{2}\\*[a-z0-9]{32}\\*[a-z0-9]{32}\\*[a-z0-9]{64}$",
		Names: []string{
			"Microsoft Office 2013",
		},
	},
	{
		Regex: "^\\$fde\\$[0-9]{2}\\$[a-f0-9]{32}\\$[0-9]{2}\\$[a-f0-9]{32}\\$[a-f0-9]{1000}[a-f0-9]{1000}[a-f0-9]{1000}[a-f0-9]{72}$",
		Names: []string{
			"Android FDE \u2264 4.3",
		},
	},
	{
		Regex: "^\\$oldoffice\\$[01]\\*[a-f0-9]{32}\\*[a-f0-9]{32}\\*[a-f0-9]{32}$",
		Names: []string{
			"Microsoft Office \u2264 2003 (MD5+RC4)",
			"Microsoft Office \u2264 2003 (MD5+RC4) collider-mode #1",
			"Microsoft Office \u2264 2003 (MD5+RC4) collider-mode #2",
		},
	},
	{
		Regex: "^\\$oldoffice\\$[34]\\*[a-f0-9]{32}\\*[a-f0-9]{32}\\*[a-f0-9]{40}$",
		Names: []string{
			"Microsoft Office \u2264 2003 (SHA1+RC4)",
			"Microsoft Office \u2264 2003 (SHA1+RC4) collider-mode #1",
			"Microsoft Office \u2264 2003 (SHA1+RC4) collider-mode #2",
		},
	},
	{
		Regex: "^(\\$radmin2\\$)?[a-f0-9]{32}$",
		Names: []string{
			"RAdmin v2.x",
		},
	},
	{
		Regex: "^{x-issha,\\s[0-9]{4}}[a-z0-9\\/+=]+$",
		Names: []string{
			"SAP CODVN H (PWDSALTEDHASH) iSSHA-1",
		},
	},
	{
		Regex: "^\\$cram_md5\\$[a-z0-9\\/+=-]+\\$[a-z0-9\\/+=-]{52}$",
		Names: []string{
			"CRAM-MD5",
		},
	},
	{
		Regex: "^[a-f0-9]{16}:2:4:[a-f0-9]{32}$",
		Names: []string{
			"SipHash",
		},
	},
	{
		Regex: "^[a-f0-9]{4,}$",
		Names: []string{
			"Cisco Type 7",
		},
	},
	{
		Regex: "^[a-z0-9\\/.]{13,}$",
		Names: []string{
			"BigCrypt",
		},
	},
	{
		Regex: "^(\\$cisco4\\$)?[a-z0-9\\/.]{43}$",
		Names: []string{
			"Cisco Type 4",
		},
	},
	{
		Regex: "^bcrypt_sha256\\$\\$(2[axy]|2)\\$[0-9]+\\$[a-z0-9\\/.]{53}$",
		Names: []string{
			"Django(bcrypt-SHA256)",
		},
	},
	{
		Regex: "^\\$postgres\\$.[^\\*]+[*:][a-f0-9]{1,32}[*:][a-f0-9]{32}$",
		Names: []string{
			"PostgreSQL Challenge-Response Authentication (MD5)",
		},
	},
	{
		Regex: "^\\$siemens-s7\\$[0-9]{1}\\$[a-f0-9]{40}\\$[a-f0-9]{40}$",
		Names: []string{
			"Siemens-S7",
		},
	},
	{
		Regex: "^(\\$pst\\$)?[a-f0-9]{8}$",
		Names: []string{
			"Microsoft Outlook PST",
		},
	},
	{
		Regex: "^sha256[:$][0-9]+[:$][a-z0-9\\/+]+[:$][a-z0-9\\/+]{32,128}$",
		Names: []string{
			"PBKDF2-HMAC-SHA256(PHP)",
		},
	},
	{
		Regex: "^(\\$dahua\\$)?[a-z0-9]{8}$",
		Names: []string{
			"Dahua",
		},
	},
	{
		Regex: "^\\$mysqlna\\$[a-f0-9]{40}[:*][a-f0-9]{40}$",
		Names: []string{
			"MySQL Challenge-Response Authentication (SHA1)",
		},
	},
	{
		Regex: "^\\$pdf\\$[24]\\*[34]\\*128\\*[0-9-]{1,5}\\*1\\*(16|32)\\*[a-f0-9]{32,64}\\*32\\*[a-f0-9]{64}\\*(8|16|32)\\*[a-f0-9]{16,64}$",
		Names: []string{
			"PDF 1.4 - 1.6 (Acrobat 5 - 8)",
		},
	},
}
