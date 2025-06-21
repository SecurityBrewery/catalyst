UI_DIR=$(dirname $0)
SCREENSHOT_DIR=../docs/screenshots/

cd ${UI_DIR}
bun screenshots

files=(
  "dashboard.png"
  "ticket.png"
  "tasks.png"
  "reactions.png"
)

for file in "${files[@]}"; do
  echo "Processing ${file}..."
  magick ${SCREENSHOT_DIR}${file} \
    \( -size 1280x720 xc:none -fill white \
    -draw "rectangle 0,0 1280,700 roundrectangle 0,700 1280,720 20,20" \) \
    -alpha on -compose DstIn -composite "${SCREENSHOT_DIR}${file}"
  magick "${SCREENSHOT_DIR}${file}" -crop 1280x719+0+0 +repage "${SCREENSHOT_DIR}${file}"
  magick composite -geometry +56+91 "${SCREENSHOT_DIR}${file}" screenshots/frame.png "${SCREENSHOT_DIR}${file}"
done
