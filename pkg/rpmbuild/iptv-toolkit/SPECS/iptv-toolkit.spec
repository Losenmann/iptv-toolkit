Name: iptv-toolkit
Version: 
Release: 1%{?dist}
Summary: 
License: 
URL: 
Source0: 

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