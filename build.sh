cd /mnt/cache/SSD/overLinks/overLinks
git pull origin main --recurse-submodules
cd frontend
git pull origin main
cd ..
cd telegram
git pull origin main
cd ..
docker build -t overlinks-be -f Dockerfile.backend .
docker build -t overlinks-fe -f Dockerfile.frontend .
docker build -t overlinks-tg -f Dockerfile.telegram .