{ pkgs ? import ../../../nix { } }:
let monchaind = (pkgs.callPackage ../../../. { });
in
monchaind.overrideAttrs (oldAttrs: {
  patches = oldAttrs.patches or [ ] ++ [
    ./broken-monchaind.patch
  ];
})
