package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/oberlademeister/hfvault"
	"github.com/urfave/cli/v2"
)

func cmdEncrypt(c *cli.Context) error {
	dataPrompt := promptui.Prompt{
		Label: "data",
	}
	data, err := dataPrompt.Run()
	if err != nil {
		return err
	}
	keyphrasePrompt := promptui.Prompt{
		Label: "keyphrase",
		Mask:  '*',
	}
	keyphrase, err := keyphrasePrompt.Run()
	if err != nil {
		return err
	}
	ciphertext, err := hfvault.EncryptToBase64([]byte(data), keyphrase)
	if err != nil {
		return nil
	}
	fmt.Println(ciphertext)
	return nil
}

func cmdDecrypt(c *cli.Context) error {
	ciphertextPrompt := promptui.Prompt{
		Label: "ciphertext",
	}
	ciphertext, err := ciphertextPrompt.Run()
	if err != nil {
		return err
	}
	keyphrasePrompt := promptui.Prompt{
		Label: "keyphrase",
		Mask:  '*',
	}
	keyphrase, err := keyphrasePrompt.Run()
	if err != nil {
		return err
	}
	cleartext, err := hfvault.DecryptFromBase64(ciphertext, keyphrase)
	if err != nil {
		return nil
	}
	fmt.Println(string(cleartext))
	return nil
}
