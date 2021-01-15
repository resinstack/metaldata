# metaldata

`metaldata` is a metadata provider for bare metal that's compatible
with the metadata needs of the [linuxkit metadata
package](https://github.com/linuxkit/linuxkit/blob/master/docs/metadata.md).

The idea is to provide the bare minimum of metadata to allow linuxkit
to successfully initialize a host on your local network.  To
accomplish this `metaldata` needs to be on the same layer 2 fabric as
the host that wishes to contact it, or you need to know what proxy ARP
is and configure it correctly (this is much harder than you think it
is).

