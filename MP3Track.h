// MP3Track.h
#pragma once
#include "Track.h"
#include <iostream>
using namespace std;

class MP3Track : public Track
{
public:
    MP3Track(string title, string artist, chrono::seconds duration, int bitrateKbps)
        : Track(std::move(title), std::move(artist), duration, "mp3"), bitrateKbps_(bitrateKbps) {}

    int bitrate() const { return bitrateKbps_; }

    void play() const override
    {
        cout << "Playing MP3: " << title_ << " by " << artist_
             << " [" << duration_.count() << "s, " << bitrateKbps_ << " kbps]" << endl;
    }

private:
    int bitrateKbps_;
};