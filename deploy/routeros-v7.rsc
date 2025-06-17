:if (([/system/resource/get total-memory] / 1024 / 1024) < 128) do={
    /log/error message="[ERROR] (IPTV Toolkit) Not enough RAM, 128MB or more required"
    :quit
}
:if (([/system/resource/get free-memory] / 1024 / 1024) < 32) do={
    /log/error message="[ERROR] (IPTV Toolkit) Not enough free RAM, need 32MB or more"
    :quit
}
:if ([/container/config/get ram-high] < 33554432) do={
    /container/config/set ram-high=33554432
    /log/info message="[INFO] (IPTV Toolkit) seated new size ram-high: 32M"
}
:if ([/container/config/get registry-url] != "https://registry-1.docker.io") do={
    /container/config/set registry-url="https://registry-1.docker.io"
    /log/info message="[INFO] (IPTV Toolkit) seated new registry-url: https://registry-1.docker.io"
}

:do {/interface/veth/add name=veth-iptv-toolkit address=172.17.0.40/24 gateway=172.17.0.1} on-error={}
:do {/interface/bridge/add name=bridge-container} on-error={}
:do {/interface/bridge/port add bridge=bridge-container interface=veth-iptv-toolkit} on-error={}
:do {/ip/address/add address=172.17.0.1/24 interface=bridge-container} on-error={}
:do {/routing/igmp-proxy/interface/add interface=bridge-container} on-error={}
/ip/firewall/nat/add chain=srcnat action=masquerade src-address=172.17.0.0/24
/ip/firewall/nat/add action=dst-nat chain=dstnat dst-address=192.168.88.1 dst-port=4022 protocol=tcp to-addresses=172.17.0.40 to-ports=4022

:do {/container/envs/add name="iptv-toolkit" key="IPTVTOOLKIT_PLAYLIST" value=""} on-error={}
:do {/container/envs/add name=iptv-toolkit key="IPTVTOOLKIT_EPG" value=""} on-error={}
:do {/container/add name="iptv-toolkit" remote-image="losenmann/iptv-toolkit:latest" interface="veth-iptv-toolkit" envlist="iptv-toolkit" start-on-boot=yes} on-error={}