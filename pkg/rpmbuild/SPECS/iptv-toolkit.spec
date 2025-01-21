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
  mkdir -p $RPM_BUILD_ROOT/%{_bindir}
  cp %{name} $RPM_BUILD_ROOT/%{_bindir}

%clean
  rm -rf $RPM_BUILD_ROOT

%files
%{_bindir}/%{name}
%changelog