if Rake::Win32.windows? then
  PYTHON='python'
  LS='dir'
  WHICH='where'
else
  PYTHON='python3'
  LS='ls -l'
  WHICH='which -a'
end

desc 'Show latest HEAD with stats'
task :what do
  sh 'go version'
  sh 'git --no-pager log -2 --stat HEAD'
end

task :tooling do
  puts "PATH is #{ENV['PATH']}"
  puts "GOPATH is #{ENV['GOPATH']}"
  puts "GOROOT is #{ENV['GOROOT']}"
  sh "#{WHICH} git || echo NA"
  sh "#{WHICH} sed || echo NA"
  sh "#{WHICH} zip || echo NA"
end

task :noassets do
  rm_f FileList['blobs/assets/micromamba.*']
  rm_f FileList['blobs/assets/*.zip']
  rm_f FileList['blobs/assets/*.yaml']
  rm_f FileList['blobs/assets/*.py']
  rm_f FileList['blobs/assets/man/*.txt']
  rm_f FileList['blobs/docs/*.md']
end

def download_link(version, platform, filename)
    "https://downloads.robocorp.com/micromamba/#{version}/#{platform}/#{filename}"
end

task :micromamba do
    version = File.read('assets/micromamba_version.txt').strip()
    puts "Using micromamba version #{version}"
    url = download_link(version, "macos64", "micromamba")
    sh "curl -o blobs/assets/micromamba.darwin_amd64 #{url}"
    url = download_link(version, "windows64", "micromamba.exe")
    sh "curl -o blobs/assets/micromamba.windows_amd64 #{url}"
    url = download_link(version, "linux64", "micromamba")
    sh "curl -o blobs/assets/micromamba.linux_amd64 #{url}"
    sh "gzip -f -9 blobs/assets/micromamba.*"
end

task :assets => [:noassets, :micromamba] do
  FileList['templates/*/'].each do |directory|
    basename = File.basename(directory)
    assetname = File.absolute_path(File.join("blobs", "assets", "#{basename}.zip"))
    rm_rf assetname
    puts "Directory #{directory} => #{assetname}"
    sh "cd #{directory} && zip -ryqD9 #{assetname} ."
  end
  cp FileList['assets/*.txt'], 'blobs/assets/'
  cp FileList['assets/*.yaml'], 'blobs/assets/'
  cp FileList['assets/*.py'], 'blobs/assets/'
  cp FileList['assets/man/*.txt'], 'blobs/assets/man/'
  cp FileList['docs/*.md'], 'blobs/docs/'
end

task :clean do
  sh 'rm -rf build/'
end

desc 'Update table of contents on docs/ directory.'
task :toc do
  sh "#{PYTHON} scripts/toc.py"
end

task :support => [:toc] do
  sh 'mkdir -p tmp build/linux64 build/macos64 build/windows64'
end

desc 'Run tests.'
task :test => [:support, :assets] do
  ENV['GOARCH'] = 'amd64'
  sh 'go test -cover -coverprofile=tmp/cover.out ./...'
  sh 'go tool cover -func=tmp/cover.out'
end

task :linux64 => [:what, :test] do
  ENV['GOOS'] = 'linux'
  ENV['GOARCH'] = 'amd64'
  sh "go build -ldflags '-s' -o build/linux64/ ./cmd/..."
  sh "sha256sum build/linux64/* || true"
end

task :macos64 => [:support] do
  ENV['GOOS'] = 'darwin'
  ENV['GOARCH'] = 'amd64'
  sh "go build -ldflags '-s' -o build/macos64/ ./cmd/..."
  sh "sha256sum build/macos64/* || true"
end

task :windows64 => [:support] do
  ENV['GOOS'] = 'windows'
  ENV['GOARCH'] = 'amd64'
  sh "go build -ldflags '-s' -o build/windows64/ ./cmd/..."
  sh "sha256sum build/windows64/* || true"
end

desc 'Setup build environment'
task :robotsetup do
    sh "#{PYTHON} -m pip install --upgrade -r robot_requirements.txt"
    sh "#{PYTHON} -m pip freeze"
end

desc 'Build local, operating system specific rcc'
task :local => [:tooling, :test] do
  sh "go build -o build/ ./cmd/..."
end

desc 'Run robot tests on local application'
task :robot => :local do
    sh "robot -L DEBUG -d tmp/output robot_tests"
end

desc 'Build commands to linux, macos, and windows.'
task :build => [:tooling, :version_txt, :linux64, :macos64, :windows64] do
  sh 'ls -l $(find build -type f)'
end

def version
  `sed -n -e '/Version/{s/^.*\`v//;s/\`$//p}' common/version.go`.strip
end

task :version_txt => :support do
  File.write('build/version.txt', "v#{version}")
end

task :default => :build
