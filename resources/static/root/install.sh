#!/usr/bin/env sh

foo() {
  local ostype = "$(uname -s)"
  local cputype = "$(uname -m)"
}

get_architecture() {
    local _ostype _cputype _bitness _arch _clibtype
    _ostype="$(uname -s)"
    _cputype="$(uname -m)"

    if [ "$_ostype" = Darwin ]; then
        # Darwin `uname -m` can lie due to Rosetta shenanigans. If you manage to
        # invoke a native shell binary and then a native uname binary, you can
        # get the real answer, but that's hard to ensure, so instead we use
        # `sysctl` (which doesn't lie) to check for the actual architecture.
        if [ "$_cputype" = i386 ]; then
            # Handling i386 compatibility mode in older macOS versions (<10.15)
            # running on x86_64-based Macs.
            # Starting from 10.15, macOS explicitly bans all i386 binaries from running.
            # See: <https://support.apple.com/en-us/HT208436>

            # Avoid `sysctl: unknown oid` stderr output and/or non-zero exit code.
            if sysctl hw.optional.x86_64 2> /dev/null || true | grep -q ': 1'; then
                _cputype=x86_64
            fi
        elif [ "$_cputype" = x86_64 ]; then
            # Handling x86-64 compatibility mode (a.k.a. Rosetta 2)
            # in newer macOS versions (>=11) running on arm64-based Macs.
            # Rosetta 2 is built exclusively for x86-64 and cannot run i386 binaries.

            # Avoid `sysctl: unknown oid` stderr output and/or non-zero exit code.
            if sysctl hw.optional.arm64 2> /dev/null || true | grep -q ': 1'; then
                _cputype=arm64
            fi
        fi
    fi

    if [ "$_ostype" = SunOS ]; then
        # Both Solaris and illumos presently announce as "SunOS" in "uname -s"
        # so use "uname -o" to disambiguate.  We use the full path to the
        # system uname in case the user has coreutils uname first in PATH,
        # which has historically sometimes printed the wrong value here.
        if [ "$(/usr/bin/uname -o)" = illumos ]; then
            _ostype=illumos
        fi

        # illumos systems have multi-arch userlands, and "uname -m" reports the
        # machine hardware name; e.g., "i86pc" on both 32- and 64-bit x86
        # systems.  Check for the native (widest) instruction set on the
        # running kernel:
        if [ "$_cputype" = i86pc ]; then
            _cputype="$(isainfo -n)"
        fi
    fi

    local _current_exe
    case "$_ostype" in

        Android)
            _ostype=linux-android
            ;;

        Linux)
            _current_exe=$(get_current_exe)
            _ostype=unknown-linux-$_clibtype
            _bitness=$(get_bitness "$_current_exe")
            ;;

        FreeBSD)
            _ostype=unknown-freebsd
            ;;

        NetBSD)
            _ostype=unknown-netbsd
            ;;

        DragonFly)
            _ostype=unknown-dragonfly
            ;;

        Darwin)
            _ostype=apple-darwin
            ;;

        illumos)
            _ostype=unknown-illumos
            ;;

        MINGW* | MSYS* | CYGWIN* | Windows_NT)
            _ostype=pc-windows-gnu
            ;;

        *)
            err "unrecognized OS type: $_ostype"
            exit 1
            ;;

    esac

    case "$_cputype" in

        i386 | i486 | i686 | i786 | x86)
            _cputype=i686
            ;;

        xscale | arm)
            _cputype=arm
            if [ "$_ostype" = "linux-android" ]; then
                _ostype=linux-androideabi
            fi
            ;;

        armv6l)
            _cputype=arm
            if [ "$_ostype" = "linux-android" ]; then
                _ostype=linux-androideabi
            else
                _ostype="${_ostype}eabihf"
            fi
            ;;

        armv7l | armv8l)
            _cputype=armv7
            if [ "$_ostype" = "linux-android" ]; then
                _ostype=linux-androideabi
            else
                _ostype="${_ostype}eabihf"
            fi
            ;;

        aarch64 | arm64)
            _cputype=aarch64
            ;;

        x86_64 | x86-64 | x64 | amd64)
            _cputype=x86_64
            ;;

        mips)
            _cputype=$(get_endianness "$_current_exe" mips '' el)
            ;;

        mips64)
            if [ "$_bitness" -eq 64 ]; then
                # only n64 ABI is supported for now
                _ostype="${_ostype}abi64"
                _cputype=$(get_endianness "$_current_exe" mips64 '' el)
            fi
            ;;

        ppc)
            _cputype=powerpc
            ;;

        ppc64)
            _cputype=powerpc64
            ;;

        ppc64le)
            _cputype=powerpc64le
            ;;

        s390x)
            _cputype=s390x
            ;;
        riscv64)
            _cputype=riscv64gc
            ;;
        loongarch64)
            _cputype=loongarch64
            ensure_loongarch_uapi
            ;;
        *)
            err "unknown CPU type: $_cputype"
            exit 1

    esac

    # Detect 64-bit linux with 32-bit userland
    if [ "${_ostype}" = unknown-linux-gnu ] && [ "${_bitness}" -eq 32 ]; then
        case $_cputype in
            x86_64)
                if [ -n "${RUSTUP_CPUTYPE:-}" ]; then
                    _cputype="$RUSTUP_CPUTYPE"
                else {
                    # 32-bit executable for amd64 = x32
                    if is_host_amd64_elf "$_current_exe"; then {
                        err "This host is running an x32 userland, for which no native toolchain is provided."
                        err "You will have to install multiarch compatibility with i686 or amd64."
                        err "To do so, set the RUSTUP_CPUTYPE environment variable set to i686 or amd64 and re-run this script."
                        err "You will be able to add an x32 target after installation by running \`rustup target add x86_64-unknown-linux-gnux32\`."
                         exit 1
                    }; else
                        _cputype=i686
                    fi
                }; fi
                ;;
            mips64)
                _cputype=$(get_endianness "$_current_exe" mips '' el)
                ;;
            powerpc64)
                _cputype=powerpc
                ;;
            aarch64)
                _cputype=armv7
                if [ "$_ostype" = "linux-android" ]; then
                    _ostype=linux-androideabi
                else
                    _ostype="${_ostype}eabihf"
                fi
                ;;
            riscv64gc)
                err "riscv64 with 32-bit userland unsupported"
                exit 1
                ;;
        esac
    fi

    # Detect armv7 but without the CPU features Rust needs in that build,
    # and fall back to arm.
    # See https://github.com/rust-lang/rustup.rs/issues/587.
    if [ "$_ostype" = "unknown-linux-gnueabihf" ] && [ "$_cputype" = armv7 ]; then
        if ! (ensure grep '^Features' /proc/cpuinfo | grep -E -q 'neon|simd') ; then
            # Either `/proc/cpuinfo` is malformed or unavailable, or
            # at least one processor does not have NEON (which is asimd on armv8+).
            _cputype=arm
        fi
    fi

    _arch="${_cputype}-${_ostype}"

    RETVAL="$_arch"
}