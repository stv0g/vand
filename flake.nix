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
        packages.vand = pkgs.buildGoModule {
          name = "vand";
          src = ./.;
          vendorHash = "sha256-dPNv9uFGbAk9Ul3YQ2woaifwez18O6plVDfd67grP+c=";
          buildInputs =
            with pkgs;
            [
              nodejs_18
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

        devShell = pkgs.mkShell {
          packages = with pkgs; [
            golangci-lint
            reuse
            gnumake
            rsync
          ];

          inputsFrom = [ packages.vand ];
        };

        formatter = nixpkgs.nixfmt-rfc-style;
      }
    );
}
