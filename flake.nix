# SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
# SPDX-License-Identifier: Apache-2.0
{
  description = "VAN daemon";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
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
          default = packages.backend;

          frontend = pkgs.buildNpmPackage {
                name = "vand-frontend";
                src = ./frontend;
                npmDepsHash = "sha256-5tXDAYY07aVnrCn7QIiPJrv9Fu+yG+Pyy9iysli6S/o=";
                installPhase = ''
                  mkdir $out
                  cp -r build/* $out/
                '';
              };

          backend = pkgs.buildGo124Module {
              name = "vand-backend";
              src = ./.;
              vendorHash = "sha256-85x6oOsfu2kQXaRhNKPRb5ZI8SKEsTAKD/5p0aFOM4E=";

              preBuild = ''
                rm -rf frontend/build
                cp -r ${packages.frontend}/ frontend/build/

                echo $PATH
                pkg-config --version
              '';

              postInstall = ''
                mv $out/bin/cmd $out/bin/vand
              '';

              tags = [
                "embed_frontend"
                "virtual"
              ];

              nativeBuildInputs = with pkgs; [
                pkg-config
                protobuf
                protoc-gen-go
              ];

              buildInputs = with pkgs; if stdenv.isLinux then
                    [
                      vulkan-headers
                      libxkbcommon
                      wayland
                      libGL
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
                    [ ];

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

          inputsFrom = with packages; [
            backend
            frontend
          ];
        };

        formatter = nixpkgs.nixfmt-rfc-style;
      }
    );
}
