#!/bin/bash
set -euo pipefail

# Post-install: enable services and apply immutability tweaks

systemctl enable NetworkManager.service
systemctl enable sshd.service
systemctl enable systemd-boot-update.service
systemctl enable dbus-broker.service

# Flatpak remote
flatpak remote-add --if-not-exists flathub https://flathub.org/repo/flathub.flatpakrepo

# Set default target to graphical
systemctl set-default graphical.target
