%define _build_id_links none
%define debug_package %{nil}
Name: iptv-toolkit
Version: 
Release: 1
Summary: A set of tools for working with IPTV
Summary(ru): Набор инструментов для работы с IPTV
License: Apache-2.0
URL: https://losenmann.github.io/iptv-toolkit
Source0: https://github.com/Losenmann/iptv-toolkit/archive/refs/tags/v%{version}.tar.gz
Group: Development/Tools

%description
Includes a play list converter and electronic program guide
Has a built-in UDP-to-HTTP and a file server

%description -l ru
Включает конвертер плейлистов и программы передач.
Имеет встроенный UDP-to-HTTP и файловый сервер.

%prep
%setup -q

%build

%install
  install -m755 -D %{name} %{buildroot}%{_bindir}/%{name}
  mkdir -p %{buildroot}%{_mandir} %{buildroot}/var/www/%{name}
  cp -r %{_topdir}/../man/* %{buildroot}%{_mandir}/

%files
  %{_bindir}/%{name}
  %{_mandir}/*/iptv-toolkit.1*
  %{_mandir}/ru/*/iptv-toolkit.1*

%check

%changelog
