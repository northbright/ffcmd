#!/bin/bash
# The command comes from running "go test -v" in root dir of ffcmd.
# Usage
# 1. Open a terminal.
# 2. Run "cd PATH/to/example" to go to current example dir.
# 3. Run "./gen.sh"
# 4. Check if "output.mp4" is generated successfully.

echo -ne "1\n00:00:00,000 --> 00:00:03,000\nGood Times with Maomi & Mimao" > "op.srt" && echo -ne "1\n00:00:00,000 --> 00:00:03,000\nMimao likes lying on father's bed...😂\nMusic by penguinmusic: Better Day" > "ed.srt" && echo -ne "1\n00:00:00,000 --> 00:00:05,000\nMido's tickling Mimao and he's enjoying..." > "01.srt" && ffprobe -v error -select_streams v:0 -show_entries stream=duration -of csv=s=,:p=0 "02.MOV" | awk -F. '{ print $1 }' | read sec; hh=$((sec / 3600)); mm=$((sec % 3600 / 60)); ss=$((sec % 3600 % 60)); printf -v end "%02d:%02d:%02d,000" hh mm ss; echo -ne "1\n00:00:00,000 --> $end\nMimao's playing the toy." > "02.srt" && echo -ne "1\n00:00:01,000 --> 00:00:09,000\nIt's hard to brush Maomi's teeth..." > "03.srt" && ffmpeg \
-i "op.jpg" \
-i "ed.jpg" \
-i "01.MP4" \
-i "02.MOV" \
-i "03.MOV" \
-i "penguinmusic-Better Day.mp3" \
-filter_complex " \
[0:v:0]fps=30,loop=loop=90:size=1,scale=720:960:force_original_aspect_ratio=decrease,pad=720:960:(ow-iw)/2:(oh-ih)/2,setsar=1:1,format=pix_fmts=yuv420p,subtitles='op.srt':force_style='Fontsize=15',fade=t=out:st=2:d=1[op_v];
aevalsrc=0:d=3[op_a];
[1:v:0]fps=30,loop=loop=90:size=1,scale=720:960:force_original_aspect_ratio=decrease,pad=720:960:(ow-iw)/2:(oh-ih)/2,setsar=1:1,format=pix_fmts=yuv420p,subtitles='ed.srt':force_style='Fontsize=13',fade=t=out:st=2:d=1[ed_v];
aevalsrc=0:d=3[ed_a];
[2:v:0]scale=720:960:force_original_aspect_ratio=decrease,pad=720:960:(ow-iw)/2:(oh-ih)/2,setsar=1:1,trim=end=5.000,setpts=PTS-STARTPTS,subtitles='01.srt':force_style='Fontsize=13'[clip_00_v];
[2:a:0]atrim=end=5.000,asetpts=PTS-STARTPTS[clip_00_a];
[3:v:0]scale=720:960:force_original_aspect_ratio=decrease,pad=720:960:(ow-iw)/2:(oh-ih)/2,setsar=1:1,subtitles='02.srt':force_style='Fontsize=13'[clip_01_v];
[4:v:0]scale=720:960:force_original_aspect_ratio=decrease,pad=720:960:(ow-iw)/2:(oh-ih)/2,setsar=1:1,trim=start=1.000:end=9.000,setpts=PTS-STARTPTS,subtitles='03.srt':force_style='Fontsize=13'[clip_02_v];
[4:a:0]atrim=start=1.000:end=9.000,asetpts=PTS-STARTPTS[clip_02_a];
[op_v][op_a][clip_00_v][clip_00_a][clip_01_v][3:a:0][clip_02_v][clip_02_a][ed_v][ed_a]concat=n=5:v=1:a=1[outv][outa];
[5:a:0][outa]amerge=inputs=2,pan=stereo|c0<c0+c2|c1<c1+c3[outa_merged_bgm]" \
-map "[outv]" \
-map "[outa_merged_bgm]" \
output.mp4 && rm "op.srt" && rm "ed.srt" && rm "01.srt" && rm "02.srt" && rm "03.srt"
