# node-icns

A command line tool to generate .icns files from a PNG file size 1024 or greater.

Obviously only meant to run on OSX.

## Install

`npm install -g node-icns`

## Usage

The file being used needs to be 1024x1024.

`nicns [--in <inFile>, --out <outFile>]`

If you do not provide an `--in <inFile>`, the current directory name is used and it will search for a PNG in the folder.

If you do not provide an `--out <outFile>`, the current directory name is used for the .icns file.

## License

node-icns is copyright 2015 typefoo, a division of uh-sem-blee, Co.

The Affero General Public License (AGPL)

Copyright (c) 2016, typefoo a division of uh-sem-blee, Co. (www.uh-sem-blee.com)

This program is free software: you can redistribute it and/or modify it under the terms of the Affero GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the Affero GNU General Public License along with this program. If not, see http://www.gnu.org/licenses/.
