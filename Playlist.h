// Playlist.h
#pragma once
#include "Track.h"
#include "MP3Track.h"
#include "WAVTrack.h"
#include <vector>
#include <memory>
#include <string>
#include <algorithm>
#include <numeric>
#include <random>
#include <chrono>
#include <iostream>
#include <fstream>
#include <filesystem>
#include <thread>

using namespace std;
namespace fs = std::filesystem;

template <typename T>
using list = vector<T>;

class Playlist
{
private:
    string name_;
    list<shared_ptr<Track>> tracks_;

    static inline const fs::path SAVE_FOLDER = fs::current_path() / "playlists";

    static void ensureFolderExists()
    {
        if (!fs::exists(SAVE_FOLDER))
        {
            if (fs::create_directories(SAVE_FOLDER))
            {
                cout << "Created playlist directory: " << SAVE_FOLDER << endl;
            }
            else
            {
                cerr << "Error: Could not create playlist directory: " << SAVE_FOLDER << endl;
                throw runtime_error("Failed to create playlist directory.");
            }
        }
    }

public:
    explicit Playlist(string name) : name_(std::move(name))
    {
        ensureFolderExists();
    }

    void addTrack(shared_ptr<Track> track)
    {
        tracks_.push_back(std::move(track));
    }

    void removeTrack(size_t index)
    {
        if (index < tracks_.size())
        {
            tracks_.erase(tracks_.begin() + index);
        }
    }

    void playAll() const
    {
        cout << "Playing playlist: " << name_ << endl;
        for (const auto &track : tracks_)
        {
            cout << "Now playing: " << track->title()
                 << " by " << track->artist()
                 << " [" << track->duration().count() << " seconds]" << endl;

            track->play();
            this_thread::sleep_for(track->duration());
            cout << "Finished: " << track->title() << "\n"
                 << endl;
        }
    }

    void shuffleTracks()
    {
        random_device rd;
        mt19937 g(rd());
        shuffle(tracks_.begin(), tracks_.end(), g);
    }

    chrono::seconds getTotalDuration() const
    {
        return accumulate(
            tracks_.begin(),
            tracks_.end(),
            chrono::seconds{0},
            [](chrono::seconds sum, const shared_ptr<Track> &track)
            {
                return sum + track->duration();
            });
    }

    void saveToDisk(const string &filenameHint = "") const
    {
        ensureFolderExists();

        string safeName = filenameHint.empty() ? name_ : filenameHint;
        fs::path filepath = SAVE_FOLDER / (safeName + ".txt");

        ofstream outFile(filepath);
        if (!outFile)
        {
            cerr << "Error: Could not open file for writing: " << filepath << endl;
            return;
        }

        outFile << name_ << "\n";
        for (const auto &track : tracks_)
        {
            string extra = "0";
            if (dynamic_cast<MP3Track *>(track.get()))
            {
                extra = to_string(static_cast<MP3Track *>(track.get())->bitrate());
            }
            else if (dynamic_cast<WAVTrack *>(track.get()))
            {
                extra = to_string(static_cast<WAVTrack *>(track.get())->sampleRate());
            }

            outFile << track->title() << "|"
                    << track->artist() << "|"
                    << track->duration().count() << "|"
                    << track->format() << "|"
                    << extra << "\n";
        }

        cout << "Playlist saved to: " << filepath << endl;
    }

    static vector<string> getSavedPlaylists()
    {
        vector<string> names;
        ensureFolderExists();

        for (const auto &entry : fs::directory_iterator(SAVE_FOLDER))
        {
            if (entry.is_regular_file() && entry.path().extension() == ".txt")
            {
                names.push_back(entry.path().stem().string());
            }
        }

        return names;
    }

    static Playlist loadFromDisk(const string &name)
    {
        ensureFolderExists();
        fs::path filepath = SAVE_FOLDER / (name + ".txt");

        ifstream inFile(filepath);
        if (!inFile)
        {
            throw runtime_error("Could not open playlist file: " + filepath.string());
        }

        string nameLine;
        getline(inFile, nameLine);
        if (nameLine.empty())
        {
            throw runtime_error("Invalid playlist file: empty name.");
        }

        Playlist playlist(nameLine);

        string line;
        while (getline(inFile, line))
        {
            if (line.empty())
                continue;

            vector<string> parts;
            size_t start = 0, pos;
            while ((pos = line.find('|', start)) != string::npos)
            {
                parts.push_back(line.substr(start, pos - start));
                start = pos + 1;
            }
            parts.push_back(line.substr(start));

            if (parts.size() != 5)
                continue;

            string title = parts[0];
            string artist = parts[1];
            int durationSec = stoi(parts[2]);
            string format = parts[3];
            int extra = stoi(parts[4]);

            shared_ptr<Track> track;
            if (format == "mp3")
            {
                track = make_shared<MP3Track>(title, artist, chrono::seconds(durationSec), extra);
            }
            else if (format == "wav")
            {
                track = make_shared<WAVTrack>(title, artist, chrono::seconds(durationSec), extra);
            }
            else
            {
                track = make_shared<Track>(title, artist, chrono::seconds(durationSec), format);
            }

            playlist.addTrack(std::move(track));
        }

        return playlist;
    }

    static void playSavedPlaylistByName(const string &name)
    {
        try
        {
            Playlist p = loadFromDisk(name);
            p.playAll();
        }
        catch (const exception &e)
        {
            cerr << "Error playing playlist '" << name << "': " << e.what() << endl;
        }
    }

    const string &name() const { return name_; }
    size_t size() const { return tracks_.size(); }

    static string getFolder() { return SAVE_FOLDER.string(); }
};