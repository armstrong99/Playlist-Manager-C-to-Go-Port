// main.cpp
#include "Playlist.h"
#include "MP3Track.h"
#include "WAVTrack.h"
#include <iostream>
#include <vector>
#include <random>
#include <thread>

using namespace std;

int main()
{
    try
    {
        auto track1 = make_shared<MP3Track>("Bohemian Rhapsody", "Queen", chrono::seconds(3), 320);
        auto track2 = make_shared<MP3Track>("Sweet Child O' Mine", "Guns N' Roses", chrono::seconds(1), 256);
        auto track3 = make_shared<WAVTrack>("Imagine", "John Lennon", chrono::seconds(1), 44100);
        auto track4 = make_shared<WAVTrack>("Billie Jean", "Michael Jackson", chrono::seconds(2), 48000);
        auto track5 = make_shared<MP3Track>("Stairway to Heaven", "Led Zeppelin", chrono::seconds(4), 192);
        auto track6 = make_shared<WAVTrack>("Smells Like Teen Spirit", "Nirvana", chrono::seconds(1), 44100);

        Playlist classicRock("Classic Rock Hits");
        classicRock.addTrack(track1);
        classicRock.addTrack(track2);
        classicRock.addTrack(track5);

        Playlist highFidelity("High-Quality WAV Collection");
        highFidelity.addTrack(track3);
        highFidelity.addTrack(track4);
        highFidelity.addTrack(track6);

        Playlist mixedPlaylist("Mixed Favorites");
        mixedPlaylist.addTrack(track1);
        mixedPlaylist.addTrack(track3);
        mixedPlaylist.addTrack(track2);
        mixedPlaylist.shuffleTracks();

        cout << "\n=== Saving playlists ===" << endl;
        classicRock.saveToDisk();
        highFidelity.saveToDisk();
        mixedPlaylist.saveToDisk();

        cout << "\n=== Discovering saved playlists ===" << endl;
        auto playlistNames = Playlist::getSavedPlaylists();

        if (playlistNames.empty())
        {
            cerr << "No playlists found on disk!" << endl;
            return 1;
        }

        for (const auto &name : playlistNames)
        {
            cout << "Found playlist: " << name << endl;
        }

        random_device rd;
        mt19937 gen(rd());
        uniform_int_distribution<size_t> dist(0, playlistNames.size() - 1);
        string selectedName = playlistNames[dist(gen)];

        cout << "\n=== Randomly selected playlist: " << selectedName << " ===" << endl;

        Playlist::playSavedPlaylistByName(selectedName);
    }
    catch (const exception &e)
    {
        cerr << "Error in main: " << e.what() << endl;
        return 1;
    }

    return 0;
}