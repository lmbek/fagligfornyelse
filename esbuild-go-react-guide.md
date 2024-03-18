# Steps

## Step 1

    npm init -y

## Step 2

    npm install --save-dev nodemon

## Step 3
Add this to the package.json file (append the start to the scripts):

    "scripts": {
        "start": "nodemon --exec 'esbuild src/index.jsx --bundle --outdir=public --target=es6 --loader:.jsx=jsx' --watch src"
    }
