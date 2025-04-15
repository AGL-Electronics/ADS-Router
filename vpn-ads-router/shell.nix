{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    (pkgs.stdenv.mkDerivation {
      pname = "go";
      version = "1.24.2";
      src = pkgs.fetchurl {
        url = "https://go.dev/dl/go1.24.2.linux-amd64.tar.gz";
        sha256 = "aAl71oCDnLydRkoO3OT3wzOXXiepAkaJDp8QeMfnAq0="; # Valid hash as of now
      };
      installPhase = ''
        mkdir -p $out
        tar -C $out -xzf $src --strip-components=1
      '';
    })
  ];
}

