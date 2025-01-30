Name: iptv-toolkit
Version: 
Release: 1%{?dist}
Summary: Toolkit for IPTV
License: 
URL: 
Source0: %{name}-%{version}.tar.gz

%description

%prep
%setup -q

%install
  rm -rf $RPM_BUILD_ROOT
  install -m755 -D %{name} %{buildroot}%{_bindir}/iptv-toolkit

%clean
  rm -rf $RPM_BUILD_ROOT

%files
  %{_bindir}/%{name}

%changelog
