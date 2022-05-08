import { yupResolver } from "@hookform/resolvers/yup";
import { Box, Button, Container, Grid, Theme } from "@mui/material";
import { createStyles, makeStyles } from '@mui/styles';
import axios from "axios";
import { FC, useState, useEffect } from "react";
import { FormProvider, SubmitHandler, useForm } from "react-hook-form";
import { SchemaOf, string, object } from "yup";

import ReactHookFormTextField from "../src/components/RHookFormText";
import { Wallet } from "../src/lib/interface/Wallet";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      flexGrow: 1,
      minHeight: '100vh',
    },
  })
);

interface IFormProps {
  sender_public_key: string
  sender_private_key: string
  recipient_blockchain_address: string
  sender_blockchain_address: string
  value: string
}

const formSchema: SchemaOf<IFormProps> = object({
  sender_public_key: string().required('入力必須です'),
  sender_private_key: string().required('入力必須です'),
  recipient_blockchain_address: string().required('入力必須です'),
  sender_blockchain_address: string().required('入力必須です'),
  value: string().required('入力必須です'),
})

const FieldArrayForm: FC = () => {
  const [wallet, setWallet] = useState<Wallet>({
    public_key: "",
    private_key: "",
    blockchain_address: "",
  })

  useEffect(() => {
    const getWallet = async () => {
      await axios.get(`${process.env.NEXT_PUBLIC_WALLET_API_URL}/wallet`)
        .then(res => {
          console.log(res.data)
          setWallet(res.data)
        })
        .catch(err => console.log(err))
    }
    getWallet()
  }, [])

  const classes = useStyles()

  const methods = useForm<IFormProps>({
    resolver: yupResolver(formSchema),
  })

  const submit: SubmitHandler<IFormProps> = async (data: IFormProps) => {
    console.log('data submitted', data)
    axios.post(`${process.env.NEXT_PUBLIC_WALLET_API_URL}/transaction`, data)
      .then((res) => console.log(res))
      .catch(err => console.log("error: ", err))
  }

  return (
    <Container maxWidth="lg">
      <Box>
        <div>
          <FormProvider {...methods}>

            <h1>Wallet</h1>
            <form onSubmit={methods.handleSubmit(submit)}>
              {wallet === undefined || wallet === null ?
                <>Loading...</>
                :
                <Grid container direction="column" sx={{ padding: 1 }}>
                  <Grid item>
                    <ReactHookFormTextField label="Public Key" name="sender_public_key" value={wallet.public_key} />
                  </Grid>

                  <Grid item>
                    <ReactHookFormTextField label="Private Key" name="sender_private_key" value={wallet.private_key} />
                  </Grid>
                  <Grid item>
                    <ReactHookFormTextField label="Blockchain Address" name="sender_blockchain_address" value={wallet.blockchain_address} />
                  </Grid>
                </Grid>
              }
              <h1>Send Money</h1>
              <Grid container direction="column" sx={{ padding: 1 }}>
                <Grid item>
                  <ReactHookFormTextField label="Sender Address" name="recipient_blockchain_address" />
                </Grid>
                <Grid item>
                  <ReactHookFormTextField label="amount" name="value" />
                </Grid>
              </Grid>
              <Grid item>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                >
                  送信
                </Button>
              </Grid>
            </form>
          </FormProvider>
        </div>
      </Box >
    </Container >
  )
}

export default FieldArrayForm
