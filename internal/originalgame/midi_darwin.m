//go:build darwin && cgo

#import <AVFoundation/AVFoundation.h>
#import <stdint.h>

static AVMIDIPlayer *dr_midi_player = nil;

int dr_play_midi(const uint8_t *bytes, int length) {
    if (bytes == NULL || length <= 0) {
        return 0;
    }
    @autoreleasepool {
        NSData *data = [NSData dataWithBytes:bytes length:(NSUInteger)length];
        NSError *error = nil;
        AVMIDIPlayer *next = [[AVMIDIPlayer alloc] initWithData:data soundBankURL:nil error:&error];
        if (next == nil) {
            return 0;
        }
        if (dr_midi_player != nil) {
            [dr_midi_player stop];
            [dr_midi_player release];
        }
        dr_midi_player = next;
        [dr_midi_player prepareToPlay];
        [dr_midi_player play:nil];
        return 1;
    }
}

void dr_stop_midi(void) {
    @autoreleasepool {
        if (dr_midi_player != nil) {
            [dr_midi_player stop];
            [dr_midi_player release];
            dr_midi_player = nil;
        }
    }
}
