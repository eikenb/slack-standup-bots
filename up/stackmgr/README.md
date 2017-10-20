Build a custom stackmgr for deployment. Customized only in that it sets the
InstanceType of the nodes.

To use run.. (see NOTE below)

    export DISABLE_STACKMGR_AUTO_UPDATE=yes
    docker build -t stackmgr .

Then you can run stackmgr normally, and it will use the custom version.

You'll need to set that environment variable in each new shell or the custom
stackmgr will get automatically overwritten by the latest cscr/stackmgr from
the docker repo.

NOTE: This assumes you're stackmgr aliases are up-to-date. If you still use the
non-custom stackmgr, update your aliases. I changed them to better support
custom stackmgrs.




