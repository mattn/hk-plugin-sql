# hk-plugin-sql

# Installation

1. Add HKPATH environment into your .bashrc or something others.

  If you are using unix, `/usr/local/lib/hk/plugin` is default.

2. Create plugins directory

        $ mkdir -p $HKPATH

3. Install sql command

        $ git clone https://github.com/mattn/hk-plugin-sql
        $ cd hk-plugin-sql
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
