#!/bin/bash
# jack-start.sh PCH -- start jack with sequencer control
# requirements: jack2 and a2j
OUTPUT=${1:-PCH} # default to internal pch sound output

jack_control start
jack_control ds alsa             # alsa driver
jack_control dps device hw:${OUTPUT}
jack_control dps rate 48000      # 48khz - lower if cracking noises
jack_control dps nperiods 2      # standard
jack_control dps period 64       # standard

sleep 5
a2j_control --ehw
a2j_control --start

# jack_lsp    -- list ports
# jack_alias system:capture_1 mainInLeft -- add alias to port
# jack_lsp -A -- list ports with aliases
