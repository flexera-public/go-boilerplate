#! /usr/bin/awk
BEGIN { FS = " " }
NF == 3 {
	k = $1 " " $2;
	if (k in block) block[k] += $3
	else block[k] = $3
}
END {
	for (k in block) {
		print k, block[k]
	}
}
