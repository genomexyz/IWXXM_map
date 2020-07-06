#!/bin/bash

# Good!
for f in RAW/*.txt; do
	ukuran=`stat --printf="%s\n" $f`
	if [ $ukuran -lt 20 ]
	then
		echo "invalid file, pass.."
	else
		f_baru=${f::12}
		echo cp $f RAW/selected/$f
		cp $f RAW/selected/
	fi
done
