# hk-plugin-sql

SQL Console for [hk](https://github.com/heroku/hk)

# Installation

1. Add HKPATH environment into your .bashrc or something others.

  If you are using unix, `/usr/local/lib/hk/plugin` is default.

2. Create plugins directory

        $ mkdir -p $HKPATH

3. Install sql command

        $ git clone https://github.com/mattn/hk-plugin-sql
        $ cd hk-plugin-sql
        $ go get -u github.com/mattn/go-runewidth
        $ go get -u github.com/lib/pq
        $ go build sql.go
        $ cp sql $HKPATH

# Usage

    $ cd /path/to/your/heroku/project
    $ hk sql
    SQL> select * from photos;
   
# License

MIT

# Author

Yasuhiro Matsumoto
