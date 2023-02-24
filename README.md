# pake

[![travis](https://travis-ci.org/schollz/pake.svg?branch=master)](https://travis-ci.org/schollz/pake) 
[![go report card](https://goreportcard.com/badge/github.com/schollz/pake)](https://goreportcard.com/report/github.com/schollz/pake)
[![Coverage Status](https://coveralls.io/repos/github/schollz/pake/badge.svg)](https://coveralls.io/github/schollz/pake)
[![godocs](https://godoc.org/github.com/schollz/pake?status.svg)](https://godoc.org/github.com/schollz/pake) 

This library will help you allow two parties to generate a mutual secret key by using a weak key that is known to both beforehand (e.g. via some other channel of communication). This is a simple API for an implementation of password-authenticated key exchange (PAKE). This protocol is derived from [Dan Boneh and Victor Shoup's cryptography book](https://crypto.stanford.edu/~dabo/cryptobook/BonehShoup_0_4.pdf) (pg 789, "PAKE2 protocol). I decided to create this library so I could use PAKE in my file-transfer utility, [croc](https://github.com/schollz/croc).


## Install

```
go get -u github.com/schollz/pake/v3
```

## Usage 

![Explanation of algorithm](https://i.imgur.com/s7oQWVP.png)

```golang
// both parties should have a weak key
weakKey := []byte{1, 2, 3}

// initialize A
A, err := pake.InitCurve(weakKey, 0, "siec")
if err != nil {
    panic(err)
}
// initialize B
B, err := pake.InitCurve(weakKey, 1, "siec")
if err != nil {
    panic(err)
}

// send A's stuff to B
err = B.Update(A.Bytes())
if err != nil {
    panic(err)
}

// send B's stuff to A
err = A.Update(B.Bytes())
if err != nil {
    panic(err)
}

// both P and Q now have strong key generated from weak key
kA, _ := A.SessionKey()
kB, _ := B.SessionKey()
fmt.Println(bytes.Equal(kA, kB))
// Output: true
```

When passing *P* and *Q* back and forth, the structure is being marshalled using `Bytes()`, which prevents any private variables from being accessed from either party.

Each function has an error. The error become non-nil when some part of the algorithm fails verification: i.e. the points are not along the elliptic curve, or if a hash from either party is not identified. If this happens, you should abort and start a new PAKE transfer as it would have been compromised. 

## Hard-coded elliptic curve points

The elliptic curve points are hard-coded to prevent an application from allowing users to supply their own points (which could be backdoors by choosing points with known discrete logs). Public points can be verified [via sage](https://sagecell.sagemath.org/?z=eJzNltuOG8cRhu8X2Hcg1hcmYWrTVdWnMuIAM8MZIfCNA-dCsCAtenq6IyLU7oakkl0Eenf_Q2olWwmcBDGQDEDOqbu6Dt9fPWm3u8nv9n8th8U3i7-_v7y4vPhi0c0P2DnSy4t7PHchqjfWeoreRaMhEIXorLIzVqyoE2XvxaqyCEc28480cAzKxmCAd95GUrUw2cCkjd57npfrcdfvdtv74zafFl4-H5b3q_VLs27WZk1r82p1eZE--vny6pN7V6_myWenv_99311efHv92-3hL_uj_A5v_vAuTfsEu8O27KblM4Gh-y1e8GvisPgKZ3anE_Hp5PG_pMWzxdnG6jd8TsAPPyzvt9e3d_u3y9XqXzhN-mq1-GLRf714fM0Y9vBaYHVO5fN50vKlWxN_HtLs_E-D-e6ZY7q8-B5PzMPGaB-jMazUtRG59l3nAwVR5No2TWM2jbdtc_bVPBga_ucHwvt_8QVUjGdfDJBx0rcKkHvqtKFBWRumxrQePJu-b3jTBHataivkBmlja6Mqxg-EMpDzpMCZ-i70KqEl77gz7UbaltrBBHFBNkOMsuFO7ED9YB0NG-tbZ2QwBhQ8nH3psGLsjW-Dscb22m20l75rIQss0FrL2s2aETVOhhY8DJFjM3hMsRvZtE1DtnV9CP3QB6cch542HcBueBiauOlFbGylo-h8A2vtoKEPvVDHvWs33sOXx3ONKELA6hvovJG2M8a6LjZuaC13YUPtRjUOzlqrLihc85FCM7QbChykh8cdTPeBVR0CZm9N55xp1ZAMzcYET-KkCyb6hgN3bE2Mba82eB025L1DXm5_LV4ahzYVQ5R24EE9Ejx0nSEbh2C0cRuDWrlO24i6djY0fdv6oQ3UK0mE6_qLAk_r8XPxnsT6c_Wy82c1Erm5BUZlMuI82xg8ew3IIXqhNQG8ODHIDKGTkmM1QhaYOhEvJEiUh9thlr1DR9j_uzZRC_Uww2yR_OCNWLbM7NR4MmCPGB1dPKI9nDKfrU7i5_UKiFTxyfsQCwHAif0Y0L1rKIpK5dP4UOqYUDvW6MaCZSWPxrkc3GQx0OQSp2iTZkdk05hTFQrBR0PG1mSmJ1W6lL24KaYkSaWEUco4Tg4uK3wZs3c0GT-anJ2MpnoZc5EshQPcMnacFXUy5EcKE1UuxBmxhhox0hfnJVlrKodgJEyReCqjSDLVJsRmsbTGDPpnOZwM2VqQrMK1UApVx1hKKKNNIZuqhTzPHsw5HdEiSi55HIGNjxLG6qi6xS8C9EzOBH3YDZ4_rJ8__iNRAOjnREm0Z6LgssGeij7kcWlBAjEHRRMhA8oFWovGmaDYHxiV9NAt9tqgNrIKWg1bFxilIKuQiREgZwO2b5SeXVTMNhAQoc_ABBHLvHft_6uVPXoYRnoO3qF9KqBiFN3jOuLbAJSDfAa8QBWTgD5wn_uIPNGZkHB8YwAS0sQhgQkIAIyGgPJGPMlTykE-0qkTFe9cRSuohmutNk55KgXiAlDTGA2sZ0LdJpMTmakWX7NDfIqNIKHTlYhcleQqhFV9oYhgQs4BioryRC9IIq4pFJaZEgtmi3F-LFIjT4SNgybNvtQSsTI0B3Uj0dWhNcYws4_wFUpJXDB84hRdjmWapGROpX6kO6UYcmIeS4TMBMocKQcqVdikKVg4OMnoOY4p6AhNQq2WikFNLSeJ6Mdc2U1AFL3BZ0nOugIfUBOTxvCRfpQxTAXCRU_h7KubtGgcKzYInqCTGqulaZw4aiIbctEJRSEZXTVjzMYkb-Aa3AmIWaeQrOA6qynJFFd_BXFAC5_EcXmxfXt_tz8u3qTDm912nJ9MpS7q9na6ub_b3h6X_fpQyvTNeHW1-vryYoHjBWbjzXXd3729GR-P5bD8MP368CbRch6_up62fyqH43K1vtptj8dduVqdZ__tzXZXFn_cvysfzM3Hcf_4k7v52Jfju_3tor_ebevx5mHZX4_pUG7q6WN0tXyxWl0_PC5XnyaVh1zuj59ZebH4Ci1_jmq-u7yod_vFn8vjukcAi095ud4ey9vD8inA-_0ceL3CyPnDHqf368V38-U_ycuXeX-X-cvV-6cA__PZdJr9I_KgExg=&lang=sage&interacts=eJyLjgUAARUAuQ==) using hashes of `croc1` and `croc2`:

```python
all_curves = {}

# Curve25519
p = 57896044618658097711785492504343953926634992332820282019728792003956564819949
A = 486662

E = EllipticCurve(GF(p),[0,A,0,1,0])
all_curves["Curve25519"] = E

# SIEC
K.<isqrt3> = QuadraticField(-3)
pi = 2^127 + 2^25 + 2^12 + 2^6 + (1 - isqrt3)/2
p = ZZ(pi.norm())

E = EllipticCurve(GF(p),[0,19]) # E: y^2 = x^3 + 19
G = E([5,12])
all_curves["SIEC"] = E

# P-521
S = 0xD09E8800291CB85396CC6717393284AAA0DA64BA
p = 0x01FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
a = 0x01FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFC
b = 0x0051953EB9618E1C9A1F929A21A0B68540EEA2DA725B99B315F3B8B489918EF109E156193951EC7E937B1652C0BD3BB1BF073573DF883D2C34F1EF451FD46B503F00
Gx= 0x00C6858E06B70404E9CD9E3ECB662395B4429C648139053FB521F828AF606B4D3DBAA14B5E77EFE75928FE1DC127A2FFA8DE3348B3C1856A429BF97E7E31C2E5BD66
Gy= 0x011839296A789A3BC0045C8A5FB42C7D1BD998F54449579B446817AFBD17273E662C97EE72995EF42640C550B9013FAD0761353C7086A272C24088BE94769FD16650
n = 0x01FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFA51868783BF2F966B7FCC0148F709A5D03BB5C9B8899C47AEBB6FB71E91386409

E = EllipticCurve(GF(p),[a,b])
all_curves["P-521"] = E

# P-256
p = 115792089210356248762697446949407573530086143415290314195533631308867097853951
r = 115792089210356248762697446949407573529996955224135760342422259061068512044369
s = 0xc49d360886e704936a6678e1139d26b7819f7e90
c = 0x7efba1662985be9403cb055c75d4f7e0ce8d84a9c5114abcaf3177680104fa0d
b = 0x5ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b
Gx = 0x6b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c296
Gy = 0x4fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f5 

E = EllipticCurve(GF(p),[-3,b])
G = E([Gx,Gy])
all_curves["P-256"] = E

# P-384
p = 39402006196394479212279040100143613805079739270465446667948293404245721771496870329047266088258938001861606973112319
r = 39402006196394479212279040100143613805079739270465446667946905279627659399113263569398956308152294913554433653942643
s = 0xa335926aa319a27a1d00896a6773a4827acdac73
c = 0x79d1e655f868f02fff48dcdee14151ddb80643c1406d0ca10dfe6fc52009540a495e8042ea5f744f6e184667cc722483
b = 0xb3312fa7e23ee7e4988e056be3f82d19181d9c6efe8141120314088f5013875ac656398d8a2ed19d2a85c8edd3ec2aef
Gx = 0xaa87ca22be8b05378eb1c71ef320ad746e1d3b628ba79b9859f741e082542a385502f25dbf55296c3a545e3872760ab7
Gy = 0x3617de4a96262c6f5d9e98bf9292dc29f8f41dbd289a147ce9da3113b5f0b8c00a60b1ce1d7e819d7a431d7c90ea0e5f

E = EllipticCurve(GF(p),[-3,b])
G = E([Gx,Gy])
all_curves["P-384"] = E


import hashlib

def find_point(E,seed=b""):
    X = int.from_bytes(hashlib.sha1(seed).digest(),"little")
    while True:
        try:
            return E.lift_x(E.base_field()(X)).xy()
        except:
            X += 1

    
for key,E in all_curves.items():
    print(f"key = {key}, P = {find_point(E,seed=b'croc2')}")
    print(f"key = {key}, P = {find_point(E,seed=b'croc1')}")
```

which returns

```
key = Curve25519, P = (793136080485469241208656611513609866400481671854, 10652526265787470154425996210700542961928029230996359640069967802965733206444)
key = Curve25519, P = (1086685267857089638167386722555472967068468061489, 37224361612322642494225506585912767208740592163316513552703274636161220046745)
key = SIEC, P = (793136080485469241208656611513609866400481671853, 18458907634222644275952014841865282643645472623913459400556233196838128612339)
key = SIEC, P = (1086685267857089638167386722555472967068468061489, 19593504966619549205903364028255899745298716108914514072669075231742699650911)
key = P-521, P = (793136080485469241208656611513609866400481671852, 4032821203812196944795502391345776760852202059010382256134592838722123385325802540879231526503456158741518531456199762365161310489884151533417829496019094620)
key = P-521, P = (1086685267857089638167386722555472967068468061489, 5010916268086655347194655708160715195931018676225831839835602465999566066450501167246678404591906342753230577187831311039273858772817427392089150297708931207)
key = P-256, P = (793136080485469241208656611513609866400481671852, 59748757929350367369315811184980635230185250460108398961713395032485227207304)
key = P-256, P = (1086685267857089638167386722555472967068468061489, 9157340230202296554417312816309453883742349874205386245733062928888341584123)
key = P-384, P = (793136080485469241208656611513609866400481671852, 7854890799382392388170852325516804266858248936799429260403044177981810983054351714387874260245230531084533936948596)
key = P-384, P = (1086685267857089638167386722555472967068468061489, 21898206562669911998235297167979083576432197282633635629145270958059347586763418294901448537278960988843108277491616)
```

which are the points used [in the code](https://github.com/schollz/pake/blob/master/pake.go#L76-L107).

## Contributing

Pull requests are welcome. Feel free to...

- Revise documentation
- Add new features
- Fix bugs
- Suggest improvements

## Thanks

Thanks [@tscholl2](https://github.com/tscholl2) for lots of implementation help, fixes, and developing the novel ["siec" curve](https://doi.org/10.1080/10586458.2017.1412371).


## License

MIT
