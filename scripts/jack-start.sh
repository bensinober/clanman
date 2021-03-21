#!/bin/bash
# jack-start.sh PCH -- start jack with sequencer control
# get card id from cat /proc/asound/cards
# requirements: jack2 and a2j
export JACK_NO_AUDIO_RESERVATION=1
OUTPUT=${1:-PCH} # default to internal pch sound output
jackd -dalsa -dhw:${OUTPUT},0 -r22100 -p256 -n2 -S -P -o2 -zs -Xseq

# jackdbus auto
# jack_control start
# jack_control ds alsa             # alsa driver
# jack_control dps device hw:${OUTPUT}
# jack_control dps rate 11025      # 48khz - lower if cracking noises
# jack_control dps nperiods 2      # standard
# jack_control dps period 256      # the lower the better

# jack_lsp    -- list ports
# jack_alias system:capture_1 mainInLeft -- add alias to port
# jack_lsp -A -- list ports with aliases

# log jack: tail -f ~/.log/jack/jackdbus.log