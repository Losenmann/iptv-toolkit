#!/usr/bin/env python
from jtv2xmltv import convert
import urllib.request
import subprocess
import gzip
import magic
import os
import sys
import re
import shutil
import argparse
import xml.dom.minidom as minidom
import schedule
import time

parser = argparse.ArgumentParser(
    prog='Python IPTV Toolkit',
    description='Toolkit Convert EPG, IPTV Playlists',
    epilog='https://github.com/Losenmann/iptv-toolkit'
)
parser.add_argument('-E', '--path-dst-epg', type=str, help='Dst path tvguide')
parser.add_argument('-P', '--path-dst-playlist', type=str, help='Dst path playlist')
parser.add_argument('-e', '--epg-url', type=str, help='EPG URL')
parser.add_argument('-p', '--playlist-url', type=str, help=' Playlist URL')
parser.add_argument('-t', '--playlist-tvg-url', type=str, help='Add EPG URL to Playlist')
parser.add_argument('-u', '--playlist-udpxy-url', type=str, help='Create playlist with built-in udpxy connection')
parser.add_argument('-H', '--healthcheck', action='store_true', help='Enable and start healthcheck server')
args = parser.parse_args()

class iptv():
    try:
        path_dst_tvguide = args.path_dst_epg if args.path_dst_epg else os.environ['PATH_DST_EPG']
    except:
        path_dst_tvguide = "/www/iptv/tvguide"

    try:
        path_dst_playlist = args.path_dst_playlist if args.path_dst_playlist else os.environ['PATH_DST_PLAYLIST']
    except:
        path_dst_playlist = "/www/iptv/playlist"

    def __init__(self):
        self.job()
        self.subprocess()
        self.scheduler()

    def buildTvguide(self, file, type):
        """ JTV TV Guide """
        if type == "application/zip":
            data = convert.convert_jtv_to_xmltv(file)
            xmltv = open('{}/tvguide.xml'.format(self.path_dst_tvguide), 'w', encoding='utf-8')
            xmltv.write(data)
            xmltv.close()
            xmltvgz = gzip.open('{}/tvguide.xml.gz'.format(self.path_dst_tvguide), 'wt', encoding='utf-8')
            xmltvgz.write(data)
            xmltvgz.close()
            shutil.copy2(file, '{}/tvguide.zip'.format(self.path_dst_tvguide))
            print("test1")

        """ XML TV Guide """
        if type == "text/xml":
            xmltv = open(file, 'r', encoding='utf-8')
            data = xmltv.read()
            xmltv.close()
            xmltvgz = gzip.open('{}/tvguide.xml.gz'.format(self.path_dst_tvguide), 'wt', encoding='utf-8')
            xmltvgz.write(data)
            xmltvgz.close()
            shutil.copy2(file, '{}/tvguide.xml'.format(self.path_dst_tvguide))
            print("test2")

        """ XML Compressed TV Guide """
        if type == "application/gzip":
            xmltvgz = gzip.open(file, 'rt', encoding='utf-8')
            data = xmltvgz.read()
            xmltvgz.close()
            xmltv = open('{}/tvguide.xml'.format(self.path_dst_tvguide), 'w', encoding='utf-8')
            xmltv.write(data)
            xmltv.close()
            shutil.copy2(file, '{}/tvguide.xml.gz'.format(self.path_dst_tvguide))
            print("test3")

    def buildPlaylist(self, file, type, tvguide=None, udpxy=None):
        r_udp = '^udp://@'
        """ XMLTV Playlist """
        if type == "text/plain":
            r_title = r',[^.].*$'
            r_id = 'tvg-name="[0-9]*"'
            r_icon = 'tvg-logo="[^"]*"'
            data = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<playlist xmlns=\"http://xspf.org/ns/0/\" version=\"1\">\n        <title>Custom</title>\n        <trackList>"
            data_end = "\n        </trackList>\n</playlist>\n"
            data_udpxy = data
            sep_1 = "                "
            sep_2 = "                        "
            with open(file, "rt") as f:
                lines = [line.rstrip() for line in f]
                f.close()

            for i in lines:
                if re.search(r_title, i):
                    id, title, icon, track = "", "", "", ""

                    m1 = re.search(r_id, i)
                    if m1:
                        id = "\n" + sep_2 + "<channel_id>" + re.sub('tvg-name=|"', '', m1.group(0)) + "</channel_id>"

                    m2 = re.search(r_title, i)
                    if m2:
                        title = "\n" + sep_2 + "<title>" + re.sub('^.*,', '', m2.group(0)) + "</title>"

                    m3 = re.search(r_icon, i)
                    if m3:
                        icon = "\n" + sep_2 + "<image>" + re.sub('tvg-logo=|"', '', m3.group(0)) + "</image>"

                    track = str("\n" + sep_1 + "<track>" + id + title + icon + "\n" + sep_2 + "<location>" + "{}" + "</location>" + "\n" + sep_1 + "</track>")

                if re.search('^[^(#)]', i):
                    data += str(track.format(i))
                    if udpxy and re.search(r_udp, i):
                        data_udpxy += str(track.format(udpxy + re.sub(r_udp, '', i)))
            data += data_end
            if udpxy:
                data_udpxy += data_end
            f = open('{}/playlist.xml'.format(self.path_dst_playlist), 'w', encoding='utf-8')
            f.write(data)
            f.close()
            if udpxy and re.search('</track>', data_udpxy):
                shutil.copy2(file, '{}/playlist.m3u'.format(self.path_dst_playlist))
                f = open('{}/playlist_udpxy.xml'.format(self.path_dst_playlist), 'w', encoding='utf-8')
                f.write(data)
                f.close()

        """ M3U/M3U8 Playlist """
        if type == "text/xml":
            data = '#EXTM3U{} cache=500 deinterlace=1'.format("" if not tvguide else ' url-tvg="' + tvguide + '" m3uautoload=1')
            data_udpxy = data
            f = open(file)
            dom = minidom.parseString(f.read())
            f.close()
            dom.normalize()

            for track in dom.getElementsByTagName("trackList"):
                for item in track.childNodes:
                    if item.nodeType == 1:
                        id, icon, title, location, location_udpxy = "", "", "", "", ""
                        for i in item.childNodes:
                            if i.nodeType == 1:
                                if i.tagName == 'channel_id':
                                    id = ',tvg-name="{}"'.format(i.firstChild.data)
                                if i.tagName == 'image':
                                    icon = ',tvg-logo="{}"'.format(i.firstChild.data)
                                if i.tagName == 'title':
                                    title = ',{}'.format(i.firstChild.data)
                                if i.tagName == 'location':
                                    location = '\n{}'.format(i.firstChild.data)
                                    if udpxy and re.search(r_udp, i.firstChild.data):
                                        location_udpxy = '\n{}'.format(udpxy + re.sub(r_udp, '', i.firstChild.data))
                        data += str("\n#EXTINF:-1" + id + icon + title + location)
                        if udpxy:
                            data_udpxy += str("\n#EXTINF:-1" + id + icon + title + location_udpxy)
            f = open('{}/playlist.m3u'.format(self.path_dst_playlist), 'w', encoding='utf-8')
            f.write(data)
            f.close()
            shutil.copy2(file, '{}/playlist.xml'.format(self.path_dst_playlist))
            if udpxy and re.search('#EXTINF', data_udpxy):
                f.open('{}/playlist_udpxy.m3u'.format(self.path_dst_playlist), 'w', encoding='utf-8')
                f.write(data_udpxy)
                f.close()

    def job(self):
        try:
            epg = args.epg_url if args.epg_url else os.environ['EPG_URL']
        except:
            epg = None

        try:
            playlist = args.playlist_url if args.playlist_url else os.environ['PLAYLIST_URL']
        except:
            playlist = None

        try:
            playlist_tvg = args.playlist_tvg_url if args.playlist_tvg_url else os.environ['PLAYLIST_TVG_URL']
        except:
            playlist_tvg = None

        try:
            playlist_udpxy = re.sub('/[^/]*$', '', args.playlist_udpxy_url if args.playlist_udpxy_url else os.environ['PLAYLIST_UDPXY_URL']) + '/udp/'
        except:
            playlist_udpxy = None

        if epg:
            try:
                with urllib.request.urlopen(epg) as resp:
                    f = open("/tmp/tvguide", 'wb')
                    f.write(resp.read())
                    f.close()
                    resp.close()
                self.buildTvguide("/tmp/tvguide", magic.from_file("/tmp/tvguide", mime=True))
            except:
                sys.exit(1)

        if playlist:
            try:
                with urllib.request.urlopen(playlist) as resp:
                    f = open("/tmp/playlist", 'wb')
                    f.write(resp.read())
                    f.close()
                    resp.close()
                self.buildPlaylist("/tmp/playlist", magic.from_file("/tmp/playlist", mime=True), playlist_tvg, playlist_udpxy)
            except:
                sys.exit(1)

    def scheduler(self):
        schedule.every().day.at("06:30").do(self.job)
        while True:
            schedule.run_pending()
            time.sleep(1)

    def subprocess(self):
        subprocess.run(["nginx"])
        subprocess.run(["udpxy", "-p", "4023", "-vl", "/proc/1/fd/1"])

if __name__ == "__main__":
    iptv()