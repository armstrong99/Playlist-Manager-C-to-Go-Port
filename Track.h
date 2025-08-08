// Track.h
#pragma once
#include <string>
#include <chrono>
#include <iostream>
using namespace std;

class Track
{
public:
    Track(string title, string artist, chrono::seconds duration, string format)
        : title_(std::move(title)), artist_(std::move(artist)),
          duration_(duration), format_(std::move(format)) {}

    virtual ~Track() = default;

    const string &title() const { return title_; }
    const string &artist() const { return artist_; }
    chrono::seconds duration() const { return duration_; }
    const string &format() const { return format_; }

    virtual void play() const
    {
        cout << "Playing: " << title_ << " by " << artist_ << endl;
    }

protected:
    string title_;
    string artist_;
    chrono::seconds duration_;
    string format_;
};