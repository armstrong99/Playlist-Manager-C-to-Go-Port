// WAVTrack.h
#pragma once
#include "Track.h"
#include <iostream>
using namespace std;

class WAVTrack : public Track
{
public:
    WAVTrack(string title, string artist, chrono::seconds duration, int sampleRateHz)
        : Track(std::move(title), std::move(artist), duration, "wav"), sampleRateHz_(sampleRateHz) {}

    int sampleRate() const { return sampleRateHz_; }

    void play() const override
    {
        cout << "Playing WAV: " << title_ << " by " << artist_
             << " [" << duration_.count() << "s, " << sampleRateHz_ << " Hz]" << endl;
    }

private:
    int sampleRateHz_;
};