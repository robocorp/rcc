if Rake::Win32.windows? then
  PYTHON='python'
  LS='dir'
else
  PYTHON='python3'
  LS='ls -l'
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
  sh "go get -u github.com/go-bindata/go-bindata/..."
  sh "which -a zip || echo NA"
  sh "which -a go-bindata || echo NA"
  sh "ls -l $HOME/go/bin"
end

task :assets do
  FileList['templates/*/'].each do |directory|
    basename = File.basename(directory)
    assetname = File.absolute_path(File.join("assets", "#{basename}.zip"))
    rm_rf assetname
    puts "Directory #{directory} => #{assetname}"
    sh "cd #{directory} && zip -ryqD9 #{assetname} ."
  end
  sh "$HOME/go/bin/go-bindata -o blobs/assets.go -pkg blobs assets/*.zip assets/man/*"
end

task :support do
  sh 'mkdir -p tmp build/linux64 build/linux32 build/macos64 build/windows64 build/windows32'
end

desc 'Run tests.'
task :test => [:support, :assets] do
  sh 'go test -cover -coverprofile=tmp/cover.out ./...'
  sh 'go tool cover -func=tmp/cover.out'
end

task :linux64 => [:what, :test] do
  ENV['GOOS'] = 'linux'
  ENV['GOARCH'] = 'amd64'
  sh "go build -ldflags '-s' -o build/linux64/ ./cmd/..."
  sh "sha256sum build/linux64/* || true"
end

task :linux32 => [:what, :test] do
  ENV['GOOS'] = 'linux'
  ENV['GOARCH'] = '386'
  sh "go build -ldflags '-s' -o build/linux32/ ./cmd/..."
  sh "sha256sum build/linux32/* || true"
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

task :windows32 => [:support] do
  ENV['GOOS'] = 'windows'
  ENV['GOARCH'] = '386'
  sh "go build -ldflags '-s' -o build/windows32/ ./cmd/..."
  sh "sha256sum build/windows32/* || true"
end

desc 'Setup build environment'
task :robotsetup do
    sh "#{PYTHON} -m pip install --upgrade -r robot_requirements.txt"
    sh "#{PYTHON} -m pip freeze"
end

desc 'Build local, operating system specific rcc'
task :local => [:tooling, :support, :assets] do
  sh "go build -o build/ ./cmd/..."
end

desc 'Run robot tests on local application'
task :robot => :local do
    sh "robot -L DEBUG -d tmp/output robot_tests"
end

desc 'Build commands to linux, macos, and windows.'
task :build => [:tooling, :version_txt, :linux64, :linux32, :macos64, :windows64, :windows32] do
  sh 'ls -l $(find build -type f)'
end

def version
  `sed -n -e '/Version/{s/^.*\`v//;s/\`$//p}' common/version.go`.strip
end

task :version_txt => :support do
  File.write('build/version.txt', "v#{version}")
end

task :default => :build

