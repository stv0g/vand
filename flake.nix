# SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
# SPDX-License-Identifier: Apache-2.0
{
  description = "VAN daemon";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs =
    {
      self,
      flake-utils,
      nixpkgs,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      rec {
        packages = {
          vand =
            let
              frontend = pkgs.buildNpmPackage {
                name = "vand-frontend";
                src = ./frontend;
                npmDepsHash = "sha256-5tXDAYY07aVnrCn7QIiPJrv9Fu+yG+Pyy9iysli6S/o=";
                installPhase = ''
                  mkdir $out
                  cp -r build/* $out/
                '';
              };
            in
            pkgs.buildGoModule {
              name = "vand";
              src = ./.;
              vendorHash = "sha256-dPNv9uFGbAk9Ul3YQ2woaifwez18O6plVDfd67grP+c=";

              preBuild = ''
                rm -rf frontend/build
                cp -r ${frontend}/ frontend/build/
              '';

              postInstall = ''
                mv $out/bin/cmd $out/bin/vand
              '';

              tags = [
                "embed_frontend"
                "virtual"
              ];

              buildInputs =
                with pkgs;
                [
                  protobuf
                  protoc-gen-go
                ]
                ++ (
                  if stdenv.isLinux then
                    [
                      vulkan-headers
                      libxkbcommon
                      wayland
                      libGL
                      pkg-config
                      xorg.libX11
                      xorg.libXcursor
                      xorg.libXfixes
                    ]
                  else if stdenv.isDarwin then
                    with darwin.apple_sdk_11_0;
                    [
                      MacOSX-SDK
                      frameworks.Foundation
                      frameworks.Metal
                      frameworks.QuartzCore
                      frameworks.AppKit
                    ]
                  else
                    [ ]
                )
                ++ (if stdenv.isLinux then with pkgs; [ libxkbcommon ] else [ ]);
              doCheck = false;
            };
        };

        devShell = pkgs.mkShell {
          packages = with pkgs; [
            golangci-lint
            reuse
            gnumake
            rsync
            nodejs_18
          ];

          inputsFrom = [ packages.vand ];
        };

        formatter = nixpkgs.nixfmt-rfc-style;
      }
    );
}
