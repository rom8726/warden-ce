curl \
--url 'smtp://127.0.0.1:1025' \
--mail-from peter@example.com \
--mail-rcpt jessica@example.com \
--upload-file - <<EOF
From: Peter <peter@example.com>
To: Jessica <jessica@example.com>
Subject: Little Secret

You’re awesome, don’t forget! ✨

EOF